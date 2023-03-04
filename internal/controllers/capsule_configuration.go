// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"

	capsulev1beta2 "github.com/clastix/capsule/api/v1beta2"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CapsuleConfiguration struct {
	client                      client.Client
	CapsuleConfigurationName    string
	DeprecatedCapsuleUserGroups []string
}

// nolint
var CapsuleUserGroups sets.String

func (c *CapsuleConfiguration) SetupWithManager(mgr ctrl.Manager) error {
	if len(c.DeprecatedCapsuleUserGroups) > 0 {
		CapsuleUserGroups = sets.NewString(c.DeprecatedCapsuleUserGroups...)

		return nil
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&capsulev1beta2.CapsuleConfiguration{}, builder.WithPredicates(predicate.NewPredicateFuncs(func(object client.Object) bool {
			return object.GetName() == c.CapsuleConfigurationName
		}))).
		Complete(c)
}

func (c *CapsuleConfiguration) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	capsuleConfig := &capsulev1beta2.CapsuleConfiguration{}

	if err := c.client.Get(ctx, types.NamespacedName{Name: request.Name}, capsuleConfig); err != nil {
		panic(err)
	}

	CapsuleUserGroups = sets.NewString(capsuleConfig.Spec.UserGroups...)

	return reconcile.Result{}, nil
}

func (c *CapsuleConfiguration) InjectClient(client client.Client) error {
	c.client = client

	return nil
}
