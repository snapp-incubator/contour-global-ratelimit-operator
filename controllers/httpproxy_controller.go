/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strings"

	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/xdserver"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// HTTPProxyReconciler reconciles a HTTPProxy object
type HTTPProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies,verbs=get;list;watch;
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/status,verbs=get;
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *HTTPProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	httproxy := &contourv1.HTTPProxy{}
	getErr := r.Get(ctx, req.NamespacedName, httproxy)
	if getErr != nil && errors.IsNotFound(getErr) {
		if isDeleted := parser.ContourLimitConfigs.Delete(req.Namespace, req.Name); isDeleted {
			xdserver.CreateNewSnapshot()
			logger.Info("object is deleted from xds server")

		}

		return ctrl.Result{}, nil
	} else if getErr != nil {
		logger.Error(getErr, "Error getting operator resource object")
		return ctrl.Result{}, getErr

	}
	if strings.ToLower(httproxy.Status.CurrentStatus) == "valid" {
		hasGlobalRateLimitPolicy, globalRateLimitPolicy, err := parser.ExtractDescriptorsFromHTTPProxy(httproxy)
		if err != nil {
			logger.Info(err.Error())
		}
		if hasGlobalRateLimitPolicy {
			logger.Info("successfully added to the xds server")
			if addErr := parser.ContourLimitConfigs.AddToConfig(globalRateLimitPolicy); addErr != nil {
				logger.Info(addErr.Error())
			}
			xdserver.CreateNewSnapshot()
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&contourv1.HTTPProxy{}).
		//WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
