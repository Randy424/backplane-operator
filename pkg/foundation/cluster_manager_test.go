// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package foundation

import (
	"os"
	"testing"

	v1 "github.com/stolostron/backplane-operator/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestClusterManager(t *testing.T) {
	tests := []struct {
		name                      string
		mce                       *v1.MultiClusterEngine
		imageOverrides            map[string]string
		expectedNodeSelector      map[string]string
		expectedRegistrationImage string
		expectedWorkImage         string
		expectedPlacementImage    string
	}{
		{
			name: "craete cluster manager without nodeSelector",
			mce:  &v1.MultiClusterEngine{},
			imageOverrides: map[string]string{
				"registration": "quay.io/stolostron/registration@sha256:fe95bca419976ca8ffe608bc66afcead6ef333b863f22be55df57c89ded75dda",
				"work":         "quay.io/stolostron/work@sha256:856d2151423f020952d9b9253676c1c4d462fab6722c8af4885fe2b19ccd1be0",
				"placement":    "quay.io/stolostron/placement@sha256:8d69eb89ee008bf95c2b877887e66cc1541c2407c9d7339fff8a9a973200660f",
			},
			expectedRegistrationImage: "quay.io/stolostron/registration@sha256:fe95bca419976ca8ffe608bc66afcead6ef333b863f22be55df57c89ded75dda",
			expectedWorkImage:         "quay.io/stolostron/work@sha256:856d2151423f020952d9b9253676c1c4d462fab6722c8af4885fe2b19ccd1be0",
			expectedPlacementImage:    "quay.io/stolostron/placement@sha256:8d69eb89ee008bf95c2b877887e66cc1541c2407c9d7339fff8a9a973200660f",
		},
		{
			name: "craete cluster manager with nodeSelector",
			mce:  &v1.MultiClusterEngine{Spec: v1.MultiClusterEngineSpec{NodeSelector: map[string]string{"node-role.kubernetes.io/infra": ""}}},
			imageOverrides: map[string]string{
				"registration": "quay.io/stolostron/registration@sha256:fe95bca419976ca8ffe608bc66afcead6ef333b863f22be55df57c89ded75dda",
				"work":         "quay.io/stolostron/work@sha256:856d2151423f020952d9b9253676c1c4d462fab6722c8af4885fe2b19ccd1be0",
				"placement":    "quay.io/stolostron/placement@sha256:8d69eb89ee008bf95c2b877887e66cc1541c2407c9d7339fff8a9a973200660f",
			},
			expectedNodeSelector:      map[string]string{"node-role.kubernetes.io/infra": ""},
			expectedRegistrationImage: "quay.io/stolostron/registration@sha256:fe95bca419976ca8ffe608bc66afcead6ef333b863f22be55df57c89ded75dda",
			expectedWorkImage:         "quay.io/stolostron/work@sha256:856d2151423f020952d9b9253676c1c4d462fab6722c8af4885fe2b19ccd1be0",
			expectedPlacementImage:    "quay.io/stolostron/placement@sha256:8d69eb89ee008bf95c2b877887e66cc1541c2407c9d7339fff8a9a973200660f",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := ClusterManager(test.mce, test.imageOverrides)

			os.Setenv("DIRECTORY_OVERRIDE", "../../")

			_, err := GetAddons()
			if err != nil {
				t.Errorf("expected cluster manager add-ons not found")
			}

			registrationImage, found, err := unstructured.NestedString(c.Object, "spec", "registrationImagePullSpec")
			if err != nil || !found {
				t.Errorf("expected cluster manager registrationImagePullSpec not found")
			}
			if registrationImage != test.expectedRegistrationImage {
				t.Errorf("expected registrationImagePullSpec %s, got %s", registrationImage, test.expectedRegistrationImage)
			}

			workImage, found, err := unstructured.NestedString(c.Object, "spec", "workImagePullSpec")
			if err != nil || !found {
				t.Errorf("expected cluster manager workImagePullSpec not found")
			}
			if workImage != test.expectedWorkImage {
				t.Errorf("expected workImagePullSpec %s, got %s", workImage, test.expectedWorkImage)
			}

			placementImage, found, err := unstructured.NestedString(c.Object, "spec", "placementImagePullSpec")
			if err != nil || !found {
				t.Errorf("expected cluster manager placementImagePullSpec not found")
			}
			if placementImage != test.expectedPlacementImage {
				t.Errorf("expected placementImagePullSpec %s, got %s", placementImage, test.expectedPlacementImage)
			}

			nodeSelector, found, err := unstructured.NestedMap(c.Object, "spec", "nodePlacement", "nodeSelector")
			if len(test.expectedNodeSelector) != 0 && (err != nil || !found) {
				t.Errorf("expected cluster manager NodeSelector not found")
			}

			for k, v := range test.expectedNodeSelector {
				if nodeSelector[k] != v {
					t.Errorf("expected NodeSelector %s, got %s", nodeSelector, test.expectedNodeSelector)
				}
			}
		})
	}
}
