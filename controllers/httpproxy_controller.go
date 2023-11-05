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
	"fmt"

	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/parser"
	"github.com/snapp-incubator/contour-global-ratelimit-operator/internal/xdserver"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var snapShotVersion int

// HTTPProxyReconciler reconciles a HTTPProxy object
type HTTPProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies,verbs=get;list;watch;update;patch;
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=projectcontour.io,resources=httpproxies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HTTPProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

func (r *HTTPProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	loger := log.FromContext(ctx)

	httproxy := &contourv1.HTTPProxy{}
	getErr := r.Get(ctx, req.NamespacedName, httproxy)
	if getErr != nil && errors.IsNotFound(getErr) {
		return ctrl.Result{}, nil
	} else if getErr != nil {
		loger.Error(getErr, "Error getting operator resource object")
		return ctrl.Result{}, getErr

	}
	//Todo: check if httpproxy status is valid or not
	has, globalRateLimitPolicy, err := parser.ExtractDescriptorsFromHTTPProxy(httproxy)
	if err != nil {
		loger.Info(err.Error())
	}
	if has {
		loger.Info(fmt.Sprintf("successfully added to the xds server. snapShotVersion: %v", snapShotVersion))
		parser.ContourLimitConfigs.AddToConfig(globalRateLimitPolicy)
		xdserver.CreateNewSnapshot(fmt.Sprint(snapShotVersion))
		snapShotVersion++
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&contourv1.HTTPProxy{}).
		Complete(r)
}
