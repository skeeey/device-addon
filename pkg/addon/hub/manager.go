package hub

import (
	"context"
	"embed"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	"open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/addon-framework/pkg/utils"
	addonapiv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workv1 "open-cluster-management.io/api/work/v1"
)

const (
	installationNamespace = "multicluster-controlplane-agent"
	addonName             = "device-addon"
)

//go:embed manifests/templates
var fs embed.FS

func Run(ctx context.Context, kubeConfig *rest.Config) error {
	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}

	agentName := utilrand.String(5)

	agentAddon, err := addonfactory.NewAgentAddonFactory(addonName, fs, "manifests/templates").
		WithAgentRegistrationOption(&agent.RegistrationOption{
			CSRConfigurations: agent.KubeClientSignerConfigurations(addonName, agentName),
			CSRApproveCheck:   utils.DefaultCSRApprover(agentName),
			PermissionConfig:  addonRBAC(kubeConfig),
		}).
		WithInstallStrategy(agent.InstallAllStrategy(installationNamespace)).
		WithAgentHealthProber(agentHealthProber()).
		BuildTemplateAgentAddon()
	if err != nil {
		klog.Errorf("failed to build agent %v", err)
		return err
	}

	err = mgr.AddAgent(agentAddon)
	if err != nil {
		klog.Fatal(err)
	}

	err = mgr.Start(ctx)
	if err != nil {
		klog.Fatal(err)
	}
	<-ctx.Done()

	return nil
}

func agentHealthProber() *agent.HealthProber {
	return &agent.HealthProber{
		Type: agent.HealthProberTypeWork,
		WorkProber: &agent.WorkHealthProber{
			ProbeFields: []agent.ProbeField{
				{
					ResourceIdentifier: workv1.ResourceIdentifier{
						Group:     "apps",
						Resource:  "deployments",
						Name:      "device-addon-agent",
						Namespace: installationNamespace,
					},
					ProbeRules: []workv1.FeedbackRule{
						{
							Type: workv1.WellKnownStatusType,
						},
					},
				},
			},
			HealthCheck: func(identifier workv1.ResourceIdentifier, result workv1.StatusFeedbackResult) error {
				if len(result.Values) == 0 {
					return fmt.Errorf("no values are probed for deployment %s/%s", identifier.Namespace, identifier.Name)
				}
				for _, value := range result.Values {
					if value.Name != "ReadyReplicas" {
						continue
					}

					if *value.Value.Integer >= 1 {
						return nil
					}

					return fmt.Errorf("readyReplica is %d for deployement %s/%s", *value.Value.Integer, identifier.Namespace, identifier.Name)
				}
				return fmt.Errorf("readyReplica is not probed")
			},
		},
	}
}

func addonRBAC(kubeConfig *rest.Config) agent.PermissionConfigFunc {
	return func(cluster *clusterv1.ManagedCluster, addon *addonapiv1alpha1.ManagedClusterAddOn) error {
		kubeclient, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return err
		}

		groups := agent.DefaultGroups(cluster.Name, addon.Name)

		clusterRole := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					Resources: []string{"devicedatamodels"},
					APIGroups: []string{"edge.open-cluster-management.io"},
				},
			},
		}
		_, err = kubeclient.RbacV1().ClusterRoles().Get(context.TODO(), clusterRole.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
				Namespace: cluster.Name,
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					Resources: []string{"configmaps"},
					APIGroups: []string{""},
				},
				{
					Verbs:     []string{"get", "list", "watch"},
					Resources: []string{"managedclusteraddons"},
					APIGroups: []string{"addon.open-cluster-management.io"},
				},
				{
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
					Resources: []string{"devices"},
					APIGroups: []string{"edge.open-cluster-management.io"},
				},
				{
					Verbs:     []string{"update", "patch"},
					Resources: []string{"devices/status"},
					APIGroups: []string{"edge.open-cluster-management.io"},
				},
			},
		}
		_, err = kubeclient.RbacV1().Roles(cluster.Name).Get(context.TODO(), role.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().Roles(cluster.Name).Create(context.TODO(), role, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
			},
			Subjects: []rbacv1.Subject{
				{Kind: "Group", APIGroup: "rbac.authorization.k8s.io", Name: groups[0]},
			},
		}
		_, err = kubeclient.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBinding.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		binding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
				Namespace: cluster.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
			},
			Subjects: []rbacv1.Subject{
				{Kind: "Group", APIGroup: "rbac.authorization.k8s.io", Name: groups[0]},
			},
		}

		_, err = kubeclient.RbacV1().RoleBindings(cluster.Name).Get(context.TODO(), binding.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().RoleBindings(cluster.Name).Create(context.TODO(), binding, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		return nil
	}
}
