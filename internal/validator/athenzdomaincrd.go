package validator

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (v *Validator) EnsureAthenzDomainCRDExists() error {
	targetVersion := v.c.Athenz.DomainCrdVersion
	if targetVersion == "" {
		return fmt.Errorf("validation failed: athenz.domainCrdVersion is empty in config")
	}

	// CRD itself:
	crd := &unstructured.Unstructured{}
	crd.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	})

	crdName := "athenzdomains.athenz.io"
	if err := v.k.Get(context.Background(), client.ObjectKey{Name: crdName}, crd); err != nil {
		return fmt.Errorf("🚨 CRD '%s' not found in cluster: %v", crdName, err)
	}

	// Get versions from obtained CRD spec:
	// versions := [
	// 	{
	// 		"name": "v1",
	// 		"schema": {
	// 			"openAPIV3Schema": {
	// 				"type": "object",
	// 				"x-kubernetes-preserve-unknown-fields": true
	// 			}
	// 		},
	// 		"served": true,
	// 		"storage": true
	// 	}
	// ]
	versions, found, err := unstructured.NestedSlice(crd.Object, "spec", "versions")
	if err != nil || !found {
		return fmt.Errorf("invalid CRD format: spec.versions not found")
	}

	// Check if the desired version (targetVersion) exists in the list:
	for _, item := range versions {
		if vMap, ok := item.(map[string]interface{}); ok {
			if name, _ := vMap["name"].(string); name == targetVersion {
				return nil // such version found
			}
		}
	}

	return fmt.Errorf("🚨 CRD exists, but version '%s' is not supported in spec.versions", targetVersion)
}
