// Copyright 2021 IBM Corp
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"k8s.io/klog/v2"
	"os"

	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/ppc64le-cloud/powervs-utils"

	"github.com/ppc64le-cloud/pvsadm/pkg"
	"github.com/ppc64le-cloud/pvsadm/pkg/client/events"
	"github.com/ppc64le-cloud/pvsadm/pkg/client/image"
	"github.com/ppc64le-cloud/pvsadm/pkg/client/instance"
	"github.com/ppc64le-cloud/pvsadm/pkg/client/network"
	"github.com/ppc64le-cloud/pvsadm/pkg/client/volume"
)

type PVMClient struct {
	InstanceName string
	InstanceID   string
	Region       string
	Zone         string

	PISession      *ibmpisession.IBMPISession
	InstanceClient *instance.Client
	ImgClient      *image.Client
	VolumeClient   *volume.Client
	NetworkClient  *network.Client
	EventsClient   *events.Client
}

func NewPVMClient(c *Client, instanceID, instanceName, ep string) (*PVMClient, error) {
	pvmclient := &PVMClient{}
	if instanceID == "" {
		svcs, err := c.ResourceClientV2.ListInstances(controllerv2.ServiceInstanceQuery{
			Type: "service_instance",
		})
		if err != nil {
			return pvmclient, fmt.Errorf("failed to list the resource instances: %v", err)
		}
		found := false
		for _, svc := range svcs {
			klog.V(4).Infof("Service ID: %s, region_id: %s, Name: %s", svc.Guid, svc.RegionID, svc.Name)
			klog.V(4).Infof("crn: %v", svc.Crn)
			if svc.Name == instanceName {
				instanceID = svc.Guid
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("instance: %s not found", instanceName)
		}
	}

	pvmclient.InstanceID = instanceID
	svc, err := c.ResourceClientV2.GetInstance(instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get a service with ID: %s, err: %v", instanceID, err)
	}

	pvmclient.InstanceName = svc.Name
	pvmclient.Zone = svc.RegionID
	pvmclient.Region, err = utils.GetRegion(pvmclient.Zone)
	if err != nil {
		return nil, err
	}

	os.Setenv("IBMCLOUD_POWER_API_ENDPOINT", fmt.Sprintf("%s.%s", pvmclient.Region, ep))

	pvmclient.PISession, err = ibmpisession.New(c.Config.IAMAccessToken, pvmclient.Region, pkg.Options.Debug, c.User.Account, pvmclient.Zone)
	if err != nil {
		return nil, err
	}

	pvmclient.ImgClient = image.NewClient(pvmclient.PISession, instanceID)
	pvmclient.VolumeClient = volume.NewClient(pvmclient.PISession, instanceID)
	pvmclient.InstanceClient = instance.NewClient(pvmclient.PISession, instanceID)
	pvmclient.NetworkClient = network.NewClient(pvmclient.PISession, instanceID)
	pvmclient.EventsClient = events.NewClient(pvmclient.PISession, instanceID)
	return pvmclient, nil
}
