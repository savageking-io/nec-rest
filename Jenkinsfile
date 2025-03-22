pipeline {
    agent { label 'Dckr' }

    environment {
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        DOCKER_REGISTRY = 'registry.hub.docker.com'
        DOCKER_USERNAME = 'savageking'
        APP_VERSION = '0.0.0'
        IS_TAG = false
    }

    stages {
        stage('Configure Image') {
            steps {
                script {
                    env.IS_TAG = env.GIT_BRANCH =~ /^tags\//
                    def branchName = env.GIT_BRANCH?.replaceAll('^refs/heads/', '')?.replaceAll('^tags/', '')

                    if (env.IS_TAG) {
                        env.DOCKER_IMAGE = "${DOCKER_USERNAME}/necrest"
                        echo "Building Tag. Docker Image: necrest"
                    } else if (branchName == 'main') {
                        env.DOCKER_IMAGE = "${DOCKER_USERNAME}/necrest-rc"
                        echo "Building Main. Docker Image: necrest-rc"
                    } else if (branchName == 'dev') {
                        env.DOCKER_IMAGE = "${DOCKER_USERNAME}/necrest-dev"
                        echo "Building Dev. Docker Image: necrest-dev"
                    } else {
                        error "Unsupported branch or tag: ${branchName}"
                    }
                    env.DOCKER_TAG = "experimental"
                }
            }
        }

        stage('Configure Tag') {
            steps {
                script {
                    env.APP_VERSION = readFile('VERSION').trim()

                    if (env.IS_TAG) {
                        env.DOCKER_TAG = env.APP_VERSION
                        echo "Building Tag. Version: ${enc.DOCKER_TAG}"
                    } else {
                        env.DOCKER_TAG = "${env.APP_VERSION}-${env.BUILD_NUMBER}"
                        echo "Building Branch. Version: ${enc.DOCKER_TAG}"
                    }

                    echo "Docker Image: ${env.DOCKER_IMAGE}:${env.DOCKER_TAG}"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    echo "Building Image: ${env.DOCKER_IMAGE}:${env.DOCKER_TAG}"
                    //def dockerImage = docker.build("${env.DOCKER_IMAGE}:${env.DOCKER_TAG}")
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                script {
                    echo "Pushing Image: ${env.DOCKER_IMAGE}:${env.DOCKER_TAG} to ${DOCKER_REGISTRY}"
                    /*
                    docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                        def dockerImage = docker.image("${env.DOCKER_IMAGE}:${env.DOCKER_TAG}")
                        dockerImage.push()
                    }
                    */
                }
            }
        }
    }

    post {
        always {
            echo "Debug: skipping clean"
            //cleanWs()
        }
        success {
            echo 'Docker image built and pushed successfully!'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}