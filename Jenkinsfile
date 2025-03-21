pipeline {
    agent { label 'Docker Build' }

    environment {
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        DOCKER_REGISTRY = 'docker.io'
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

                    sh 'cat VERSION || type VERSION || echo "VERSION file missing"'

                    def version = readFile('VERSION').trim()
                    if (isTag) {
                        env.DOCKER_TAG = version
                        echo "Building Tag. Version: ${enc.DOCKER_TAG}"
                    } else {
                        env.DOCKER_TAG = "${version}-${env.BUILD_NUMBER}"
                        echo "Building Branch. Version: ${enc.DOCKER_TAG}"
                    }

                    echo "Docker Image: ${env.DOCKER_IMAGE}:${env.DOCKER_TAG}"
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    def dockerImage = docker.build("${env.DOCKER_IMAGE}:${env.DOCKER_TAG}")
                }
            }
        }

        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
                        def dockerImage = docker.image("${env.DOCKER_IMAGE}:${env.DOCKER_TAG}")
                        dockerImage.push()
                    }
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Docker image built and pushed successfully!'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}