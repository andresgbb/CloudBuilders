pipeline {
  agent any

  environment {
    GIT_CREDENTIALS = 'github-credentials-id'
    AWS_CREDENTIALS = 'aws-creds'
    ECR_URI = '706863001550.dkr.ecr.eu-north-1.amazonaws.com/cloudbuilders'
    KUBE_CREDENTIALS = 'kubeconfig-credentials-id'
    REGION = 'eu-north-1'
  }

  parameters {
    choice(name: 'TARGET', choices: ['minikube','ecr'], description: 'Destino del build: minikube (local) o ecr (push a registry)')
    booleanParam(name: 'DEPLOY_TO_K8S', defaultValue: true, description: 'Si es true, actualiza el deployment en Kubernetes tras build')
  }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    timeout(time: 30, unit: 'MINUTES')
  }

  stages {
    stage('Checkout') {
      steps {
        checkout([
          $class: 'GitSCM',
          branches: [[name: '*/main']],
          userRemoteConfigs: [[
            url: 'https://github.com/andresgbb/CloudBuilders.git',
            credentialsId: env.GIT_CREDENTIALS
          ]]
        ])
      }
    }

    stage('Set Image Tag') {
      steps {
        script {
          env.GIT_COMMIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
          env.IMAGE_TAG = "${env.GIT_COMMIT_SHORT}-${env.BUILD_NUMBER}"
          echo "Image tag: ${env.IMAGE_TAG}"
        }
      }
    }

    stage('Go Test') {
      steps {
        sh 'go version || true'
        sh 'go test ./...'
      }
    }

    stage('Build image') {
      steps {
        script {
          if (params.TARGET == 'minikube') {
            echo "Building image for Minikube"
            sh '''
              set -euo pipefail
              # ensure docker and minikube available
              command -v docker >/dev/null 2>&1 || { echo "docker not found"; exit 1; }
              command -v minikube >/dev/null 2>&1 || { echo "minikube not found"; exit 1; }
              # use minikube's docker daemon so image is available to cluster
              eval $(minikube -p minikube docker-env)
              docker build -t cloudbuilders:${IMAGE_TAG} .
              # optionally tag latest
              docker tag cloudbuilders:${IMAGE_TAG} cloudbuilders:latest || true
            '''
          } else {
            echo "Building and pushing image to ECR"
            withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: env.AWS_CREDENTIALS]]) {
              sh '''
                set -euo pipefail
                command -v docker >/dev/null 2>&1 || { echo "docker not found"; exit 1; }
                # Login to ECR
                aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ECR_URI%/*}
                # Build and push
                docker build -t ${ECR_URI}:${IMAGE_TAG} .
                docker push ${ECR_URI}:${IMAGE_TAG}
                docker tag ${ECR_URI}:${IMAGE_TAG} ${ECR_URI}:latest || true
                docker push ${ECR_URI}:latest || true
              '''
            }
          }
        }
      }
    }

    stage('Prepare imagePullSecret (ECR only)') {
      when { expression { return params.TARGET == 'ecr' && params.DEPLOY_TO_K8S } }
      steps {
        withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: env.AWS_CREDENTIALS]]) {
          sh '''
            set -euo pipefail
            # create or replace regcred in default namespace using current kubeconfig
            PASSWORD=$(aws ecr get-login-password --region ${REGION})
            kubectl create secret docker-registry regcred \
              --docker-server=${ECR_URI%/*} \
              --docker-username=AWS \
              --docker-password="${PASSWORD}" \
              --dry-run=client -o yaml | kubectl apply -f -
          '''
        }
      }
    }

    stage('Deploy to Kubernetes') {
      when { expression { return params.DEPLOY_TO_K8S } }
      steps {
        withCredentials([file(credentialsId: env.KUBE_CREDENTIALS, variable: 'KUBECONFIG_FILE')]) {
          sh '''
            set -euo pipefail
            export KUBECONFIG=${KUBECONFIG_FILE}
            if [ "${TARGET}" = "minikube" ]; then
              # use local image name available in minikube's docker daemon
              kubectl set image deployment/cloudbuilders app=cloudbuilders:${IMAGE_TAG} --record
            else
              kubectl set image deployment/cloudbuilders app=${ECR_URI}:${IMAGE_TAG} --record
            fi
            kubectl rollout status deployment/cloudbuilders --timeout=180s
          '''
        }
      }
    }

    stage('Cleanup') {
      steps {
        sh 'docker image prune -f || true'
      }
    }
  }

  post {
    success {
      echo "Build ${env.BUILD_NUMBER} succeeded. Image tag: ${IMAGE_TAG}"
    }
    failure {
      echo "Build ${env.BUILD_NUMBER} failed."
    }
    always {
      archiveArtifacts artifacts: '**/logs/**/*.log', allowEmptyArchive: true
    }
  }
}