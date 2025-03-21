pipeline {
    agent { label 'Docker Build' }

    environment {
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        DOCKER_REGISTRY = 'registry.hub.docker.com'
        DOCKER_USERNAME = 'savageking'
    }

    stages {
        stage('Configure Build') {
            steps {
                script {
                    def isTag = env.GIT_BRANCH =~ /^tags\//
                    def branchName = env.GIT_BRANCH?.replaceAll('^refs/heads/', '')?.replaceAll('^tags/', '')

                    if (isTag) {
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

                    def app_version
                    try {
                        app_version = sh(script: 'cat VERSION || type VERSION', returnStdout: true).trim()
                        echo "VERSION content: '${app_version}'"
                        if (app_version == '') {
                            error "VERSION file is empty or could not be read."
                        }
                    } catch (Exception e) {
                        error "Failed to get version from VERSION file: ${e.message}. Ensure VERSION exists in repo root."
                    }

                    if (isTag) {
                        env.DOCKER_TAG = app_version
                        echo "Building Tag. Version: ${enc.DOCKER_TAG}"
                    } else {
                        env.DOCKER_TAG = "${app_version}-${env.BUILD_NUMBER}"
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
                    //docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                    //    def dockerImage = docker.image("${env.DOCKER_IMAGE}:${env.DOCKER_TAG}")
                    //    dockerImage.push()
                    //}
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