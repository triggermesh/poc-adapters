export TERAFORMTARGET_TERRAFORM_CONFIG=
export TERAFORMTARGET_WORKING_DIR=/Users/jeffreynaef/go/src/github.com/triggermesh/poc-adapters/teraformtarget/cmd
export TERAFORMTARGET_TERRAFORM_PLAN='
provider "kubernetes" {
  #Context to choose from the config file.
  config_path    = "/project/cmd/"
  config_context = "arn:aws:eks:us-west-2:043455440429:cluster/tmkongdemo"
}
#-----------------------------------------
# KUBERNETES DEPLOYMENT COLOR APP
#-----------------------------------------
resource "kubernetes_deployment" "color" {
    metadata {
        name = "color-blue-dep"
        labels = {
            app   = "color"
            color = "blue"
        } //labels
    } //metadata

    spec {
        selector {
            match_labels = {
                app   = "color"
                color = "blue"
            } //match_labels
        } //selector
        #Template for the creation of the pod
        template {
            metadata {
                labels = {
                    app   = "color"
                    color = "blue"
                } //labels
            } //metadata
            spec {
                container {
                    image = "itwonderlab/color"   #Docker image name
                    name  = "color-blue"          #Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL).

                    #Block of string name and value pairs to set in the containers environment
                    env {
                        name = "COLOR"
                        value = "blue"
                    } //env

                    #List of ports to expose from the container.
                    port {
                        container_port = 8080
                    }//port
                } //container
            } //spec
        } //template
    } //spec
} //resource
'
