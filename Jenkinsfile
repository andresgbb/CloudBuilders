pipeline {
    agent any

    environment {
        IMAGE_NAME = 'cloudbuilders-image'
        CONTAINER_NAME = 'cloudbuilders-app'
        HOST_PORT = '8081'       // Puerto del host EC2
        CONTAINER_PORT = '8080'  // Puerto interno del contenedor
    }

    stages {
        stage('Checkout') {
            steps {
                git(url: 'https://github.com/andresgbb/CloudBuilders.git', branch: 'main')
            }
        }

        stage('Build Docker image') {
            steps {
                // Construye la imagen Docker en el host
                sh "docker build -t ${IMAGE_NAME} ."
            }
        }

        stage('Run container') {
            steps {
                // Elimina contenedores antiguos con el mismo nombre si existen
                sh """
                if [ \$(docker ps -a -q -f name=${CONTAINER_NAME}) ]; then
                    docker rm -f ${CONTAINER_NAME}
                fi
                """

                // Levanta el contenedor exponiendo HOST_PORT en el EC2
                sh "docker run -d --name ${CONTAINER_NAME} -p ${HOST_PORT}:${CONTAINER_PORT} ${IMAGE_NAME}"
            }
        }
    }

    post {
        success {
            echo "Pipeline completada correctamente. Accede a la app en http://<IP-PUBLICA-EC2>:${HOST_PORT}"
        }
        failure {
            echo 'La build falló.'
        }
    }
}
