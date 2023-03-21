package utils

import (
	"context"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"

	edgev1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
)

type UpdateDeviceStatusFunc func(status *edgev1alpha1.DeviceStatus) error

func UpdateDeviceStatus(
	ctx context.Context,
	client deviceclientset.Interface,
	namespace, name string,
	updateFuncs ...UpdateDeviceStatusFunc) (*edgev1alpha1.DeviceStatus, bool, error) {
	updated := false
	var updatedDeviceStatus *edgev1alpha1.DeviceStatus

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		device, err := client.EdgeV1alpha1().Devices(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		oldStatus := &device.Status

		newStatus := oldStatus.DeepCopy()
		for _, update := range updateFuncs {
			if err := update(newStatus); err != nil {
				return err
			}
		}
		if equality.Semantic.DeepEqual(oldStatus, newStatus) {
			// We return the newStatus which is a deep copy of oldStatus but with all update funcs applied.
			updatedDeviceStatus = newStatus
			return nil
		}

		oldData, err := json.Marshal(edgev1alpha1.Device{
			Status: *oldStatus,
		})

		if err != nil {
			return fmt.Errorf("failed to Marshal old data for cluster status %s: %w", device.Name, err)
		}

		newData, err := json.Marshal(edgev1alpha1.Device{
			ObjectMeta: metav1.ObjectMeta{
				UID:             device.UID,
				ResourceVersion: device.ResourceVersion,
			}, // to ensure they appear in the patch as preconditions
			Status: *newStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to Marshal new data for cluster status %s: %w", device.Name, err)
		}

		patchBytes, err := jsonpatch.CreateMergePatch(oldData, newData)
		if err != nil {
			return fmt.Errorf("failed to create patch for cluster %s: %w", device.Name, err)
		}

		updatedDevice, err := client.EdgeV1alpha1().Devices(namespace).Patch(
			ctx, device.Name, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")

		updatedDeviceStatus = &updatedDevice.Status
		updated = err == nil
		return err
	})

	return updatedDeviceStatus, updated, err
}

func UpdateDeviceConditionFn(cond metav1.Condition) UpdateDeviceStatusFunc {
	return func(oldStatus *edgev1alpha1.DeviceStatus) error {
		meta.SetStatusCondition(&oldStatus.Conditions, cond)
		return nil
	}
}
