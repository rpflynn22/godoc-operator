package controllers

import (
	"context"

	appsApi "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	netApi "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
	"github.com/rpflynn22/godoc-operator/internal/managed"
)

// RepoReconciler reconciles a Repo object
type RepoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *RepoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("namespace", req.Namespace, "name", req.Name)
	logger.Info("reconciling")

	repo := godocApi.Repo{}
	if err := r.Get(ctx, req.NamespacedName, &repo); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "get object")
		} else {
			logger.Info("resource not found, delete likely")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy := &appsApi.Deployment{ObjectMeta: metav1.ObjectMeta{Name: managed.ResourceName(req.Name), Namespace: req.Namespace}}
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, deploy, func() error {
		managed.UpdateDeployment(&repo, deploy)
		return controllerutil.SetControllerReference(&repo, deploy, r.Scheme)
	})
	if err != nil {
		logger.Error(err, "create/update deploy", "result", result)
		return ctrl.Result{}, err
	}

	service := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: managed.ResourceName(req.Name), Namespace: req.Namespace}}
	result, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		managed.UpdateService(&repo, service)
		return controllerutil.SetControllerReference(&repo, service, r.Scheme)
	})
	if err != nil {
		logger.Error(err, "create/update service", "result", result)
		return ctrl.Result{}, err
	}

	ingress := &netApi.Ingress{ObjectMeta: metav1.ObjectMeta{Name: managed.ResourceName(req.Name), Namespace: req.Namespace}}
	result, err = controllerutil.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		managed.UpdateIngress(&repo, ingress)
		return controllerutil.SetControllerReference(&repo, ingress, r.Scheme)
	})
	if err != nil {
		logger.Error(err, "create/update ingress", "result", result)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RepoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&godocApi.Repo{}).
		Owns(&appsApi.Deployment{}).
		Owns(&v1.Service{}).
		Owns(&netApi.Ingress{}).
		Complete(r)
}
