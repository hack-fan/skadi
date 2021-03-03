pipeline {
    agent any
    environment {
        JOB_NAME = "Skadi"
        CODING_DOCKER_REG_HOST = "${env.CCI_CURRENT_TEAM}-docker.pkg.${env.CCI_CURRENT_DOMAIN}"
        DOCKER_REPO_PREFIX = "${env.CODING_DOCKER_REG_HOST}/${env.PROJECT_NAME}/images/"
    }
    stages {
        stage('Checkout') {
            steps {
                checkout([$class: 'GitSCM', branches: [[name: env.GIT_BUILD_REF]],
                userRemoteConfigs: [[url: env.GIT_REPO_URL, credentialsId: env.CREDENTIALS_ID]]])
            }
        }
        stage('Build and Push') {
            steps {
                script {
                    docker.withRegistry("https://${env.CODING_DOCKER_REG_HOST}", "${env.CODING_ARTIFACTS_CREDENTIALS_ID}") {
                        def img1 = docker.build("${env.DOCKER_REPO_PREFIX}agent-api:latest","-f ./cmd/agentapi/Dockerfile .")
                        img1.push()
                        def img2 = docker.build("${env.DOCKER_REPO_PREFIX}watcher:latest","-f ./cmd/lonelywatcher/Dockerfile .")
                        img2.push()
                    }
                }
            }
        }
    }
}