kind create cluster
kubectl get nodes
keptn install

$tokenEncoded = $(kubectl get secret keptn-api-token -n keptn -ojsonpath='{.data.keptn-api-token}')
$Env:KEPTN_API_TOKEN = [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($tokenEncoded))
Write-Host $Env:KEPTN_API_TOKEN

Write-Host Run the following command to access your keptn installation and before authenticating the client
Write-Host kubectl -n keptn port-forward service/api-gateway-nginx 8080:80

Write-Host "kubectl -n keptn port-forward service/api-gateway-nginx 8080:80"

read-host "Press ENTER to continue..."

Write-Host API Endpoint http://localhost:8080/api
$Env:KEPTN_ENDPOINT = 'http://localhost:8080/api'

keptn auth --endpoint=$Env:KEPTN_ENDPOINT --api-token=$Env:KEPTN_API_TOKEN

Write-Host Starting a local registry
docker run -d -p 5000:5000 --restart=always --name registry registry:2
skaffold config set --global insecure-registries localhost:5000

Write-Host Configure your local docker instance to allow the local insecure registry 

keptn configure bridge --output

keptn create project dev --shipyard=shipyard.yml
keptn create service --project dev myservice

keptn add-resource --project=dev --service=myservice --stage=dev --resource=gitlab/gitlab.conf.yaml --resourceUri=gitlab/gitlab.conf.yaml

Write-Host keptn send event -f configuration-change.http