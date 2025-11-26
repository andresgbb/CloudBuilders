pipeline {
  agent any

  environment {
    GIT_CREDENTIALS = 'github-credentials-id'
    AWS_CREDENTIALS = 'aws-creds'
    ECR_URI = '706863001550.dkr.ecr.eu-north-1.amazonaws.com/cloudbuilders'
    KUBE_CREDENTIALS = 'kubeconfig-credentials-id'
    REGION = 'eu-north-1'
    K8S_NAMESPACE = 'default'   // Cambia si usas otro namespace
  }

  parameters {
    choice(name: 'TARGET', choices: ['minikube','ecr'], description: 'Destino del build: minikube (local) o ecr (push a registry)')
    booleanParam(name: 'DEPLOY_TO_K8S', defaultValue: true, description: 'Si es true, actualiza el deployment en Kubernetes tras build')
  }

  options {
    timestamps()
    buildDiscarder(logRotator(numToKeepStr: '30'))
    timeout(time: 60, unit: 'MINUTES')
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

    stage('Check Go version and Run Tests') {
      steps {
        sh '''
          set -euo pipefail
          bash -lc 'go version'
          # Requerimos Go >= 1.25
          if ! go version | grep -q "go1.25"; then
            echo "ERROR: Go >= 1.25 required on the agent. Install Go 1.25.4+ or use an agent image with Go preinstalled."
            exit 1
          fi
          go test ./...
        '''
      }
    }

    stage('Build image') {
      steps {
        script {
          if (params.TARGET == 'minikube') {
            echo "Building image for Minikube"
            sh '''
              set -euo pipefail
              bash -lc '
                command -v docker >/dev/null 2>&1 || { echo "docker not found"; exit 1; }
                command -v minikube >/dev/null 2>&1 || { echo "minikube not found"; exit 1; }
                eval $(minikube -p minikube docker-env)
                docker build -t cloudbuilders:${IMAGE_TAG} .
                docker tag cloudbuilders:${IMAGE_TAG} cloudbuilders:latest || true
              '
            '''
          } else {
            echo "Building and pushing image to ECR"
            withCredentials([[$class: 'AmazonWebServicesCredentialsBinding', credentialsId: env.AWS_CREDENTIALS]]) {
              sh '''
                set -euo pipefail
                bash -lc '
                  command -v docker >/dev/null 2>&1 || { echo "docker not found"; exit 1; }
                  command -v aws >/dev/null 2>&1 || { echo "aws cli not found"; exit 1; }
                  # calcular registry de forma segura
                  ECR_REGISTRY=$(echo "${ECR_URI}" | sed "s#/[^/]*$##")
                  # crear repo si no existe
                  aws ecr describe-repositories --repository-names cloudbuilders --region ${REGION} >/dev/null 2>&1 || \
                    aws ecr create-repository --repository-name cloudbuilders --region ${REGION}
                  # login y push (reintentos para robustez)
                  aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}
                  docker build -t ${ECR_URI}:${IMAGE_TAG} .
                  retry=0
                  until docker push ${ECR_URI}:${IMAGE_TAG}; do
                    retry=$((retry+1))
                    if [ $retry -ge 3 ]; then echo "docker push failed"; exit 1; fi
                    sleep 3
                  done
                  docker tag ${ECR_URI}:${IMAGE_TAG} ${ECR_URI}:latest || true
                  docker push ${ECR_URI}:latest || true
                '
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
            bash -lc '
              ECR_REGISTRY=$(echo "${ECR_URI}" | sed "s#/[^/]*$##")
              PASSWORD=$(aws ecr get-login-password --region ${REGION})
              kubectl --namespace ${K8S_NAMESPACE} create secret docker-registry regcred \
                --docker-server=${ECR_REGISTRY} \
                --docker-username=AWS \
                --docker-password="${PASSWORD}" \
                --dry-run=client -o yaml | kubectl --namespace ${K8S_NAMESPACE} apply -f -
            '
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
            bash -lc '
              export KUBECONFIG=${KUBECONFIG_FILE}
              # obtener nombre real del contenedor en el deployment
              CONTAINER_NAME=$(kubectl -n ${K8S_NAMESPACE} get deploy cloudbuilders -o jsonpath="{.spec.template.spec.containers[0].name}")
              if [ "${TARGET}" = "minikube" ]; then
                kubectl -n ${K8S_NAMESPACE} set image deployment/cloudbuilders ${CONTAINER_NAME}=cloudbuilders:${IMAGE_TAG} --record
              else
                kubectl -n ${K8S_NAMESPACE} set image deployment/cloudbuilders ${CONTAINER_NAME}=${ECR_URI}:${IMAGE_TAG} --record
              fi
              kubectl -n ${K8S_NAMESPACE} rollout status deployment/cloudbuilders --timeout=180s
            '
          '''
        }
      }
    }

    stage('Smoke test and optional rollback') {
      when { expression { return params.DEPLOY_TO_K8S } }
      steps {
        withCredentials([file(credentialsId: env.KUBE_CREDENTIALS, variable: 'KUBECONFIG_FILE')]) {
          sh '''
            set -euo pipefail
            bash -lc '
              export KUBECONFIG=${KUBECONFIG_FILE}
              # intentar un smoke test simple contra el service (ajusta si tu service usa otro puerto/nombre)
              SERVICE_IP=$(kubectl -n ${K8S_NAMESPACE} get svc cloudbuilders-service -o jsonpath="{.spec.clusterIP}" 2>/dev/null || true)
              if [ -n "${SERVICE_IP}" ]; then
                # prueba desde dentro del cluster usando un pod temporal con curl
                kubectl -n ${K8S_NAMESPACE} run --rm -i --restart=Never smoke-test --image=appropriate/curl --command -- curl -sSf http://${SERVICE_IP}:80/ || {
                  echo "Smoke test failed, rolling back"
                  kubectl -n ${K8S_NAMESPACE} rollout undo deployment/cloudbuilders
                  exit 1
                }
              else
                echo "No cloudbuilders-service encontrado; omitiendo smoke test de servicio."
              fi
            '
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