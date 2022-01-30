package aws

import (
	"encoding/json"
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/craigjames16/cfstate/utils"
)

var (
	CFService *cf.CloudFormation
)

type AppInput struct {
	Name     string
	StackID  string
	Template string
	Config   string
}
type CreateUpdateOutput struct {
	StackID string
}

type CFAppConfig struct {
	ParameterKey     string
	ParameterValue   string
	UsePreviousValue bool
}

func init() {
	var (
		sess *session.Session
		svc  *cf.CloudFormation
	)

	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1"),
	}))

	svc = cf.New(sess)

	CFService = svc
}

func CreateStack(app AppInput) (output CreateUpdateOutput, err error) {
	var (
		templateData   []byte
		configData     []byte
		cfAppConfig    []*cf.Parameter
		cfCreateOutput *cf.CreateStackOutput
	)

	templateData, err = utils.GetFileData(app.Template)
	if err != nil {
		return output, err
	}

	configData, err = utils.GetFileData(app.Config)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(configData, &cfAppConfig)
	if err != nil {
		return output, err
	}

	now := time.Now()
	sec := now.Unix()

	var (
		params = &cf.CreateStackInput{
			StackName:    aws.String(fmt.Sprintf("%s-%d", app.Name, sec)),
			TemplateBody: aws.String(string(templateData)),
			Parameters:   cfAppConfig,
		}
	)

	cfCreateOutput, err = CFService.CreateStack(params)
	if err != nil {
		return output, err
	}

	output.StackID = *cfCreateOutput.StackId

	return output, nil
}

func UpdateStack(app AppInput) (output CreateUpdateOutput, err error) {
	var (
		templateData   []byte
		configData     []byte
		cfAppConfig    []*cf.Parameter
		cfUpdateOutput *cf.UpdateStackOutput
	)

	templateData, err = utils.GetFileData(app.Template)
	if err != nil {
		return output, err
	}

	configData, err = utils.GetFileData(app.Config)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(configData, &cfAppConfig)
	if err != nil {
		return output, err
	}

	var (
		params = &cf.UpdateStackInput{
			StackName:    aws.String(app.StackID),
			TemplateBody: aws.String(string(templateData)),
			Parameters:   cfAppConfig,
		}
	)

	cfUpdateOutput, err = CFService.UpdateStack(params)
	if err != nil {
		return output, err
	}

	output.StackID = *cfUpdateOutput.StackId

	return output, nil
}
