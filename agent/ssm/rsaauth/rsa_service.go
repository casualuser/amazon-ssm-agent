// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the AWS Customer Agreement (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/agreement/

// Package rsaauth is an interface to the RSA signed methods of the SSM service.
package rsaauth

import (
	"github.com/aws/amazon-ssm-agent/agent/ssm/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/signer/v4"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// RsaSignedService is an interface to the RSA signed methods of the SSM service.
type RsaSignedService interface {
	RequestManagedInstanceRoleToken(fingerprint string) (response *ssm.RequestManagedInstanceRoleTokenOutput, err error)
	UpdateManagedInstancePublicKey(publicKey, publicKeyType string) (response *ssm.UpdateManagedInstancePublicKeyOutput, err error)
}

// sdkService is an service wrapper that delegates to the ssm sdk.
type sdkService struct {
	sdk *ssm.SSM
}

// NewService creates a new SSM service instance.
func NewRsaService(serverId string, region string, encodedPrivateKey string) RsaSignedService {
	awsConfig := util.AwsConfig()

	awsConfig.Region = &region
	awsConfig.Credentials = credentials.NewStaticCredentials(serverId, encodedPrivateKey, "")

	// Create a session to share service client config and handlers with
	ssmSess := session.New(awsConfig)

	// Clear existing singers
	ssmSess.Handlers.Sign.Clear()
	// Add custom signer to session, will be used by any service created with this session
	ssmSess.Handlers.Sign.PushBack(v4.SignRsa)

	ssmService := ssm.New(ssmSess)
	return &sdkService{sdk: ssmService}
}

// RequestManagedInstanceRoleToken calls the RequestManagedInstanceRoleToken SSM API.
func (svc *sdkService) RequestManagedInstanceRoleToken(fingerprint string) (response *ssm.RequestManagedInstanceRoleTokenOutput, err error) {

	params := ssm.RequestManagedInstanceRoleTokenInput{
		Fingerprint: aws.String(fingerprint),
	}

	return svc.sdk.RequestManagedInstanceRoleToken(&params)
}

// UpdateManagedInstancePublicKey calls the UpdateManagedInstancePublicKey SSM API.
func (svc *sdkService) UpdateManagedInstancePublicKey(publicKey, publicKeyType string) (response *ssm.UpdateManagedInstancePublicKeyOutput, err error) {

	params := ssm.UpdateManagedInstancePublicKeyInput{
		NewPublicKey:     aws.String(publicKey),
		NewPublicKeyType: aws.String(publicKeyType),
	}

	return svc.sdk.UpdateManagedInstancePublicKey(&params)
}
