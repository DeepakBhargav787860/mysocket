
 //https://mysocket-6xmu.onrender.com
 //socket-application-mp9v-hlpxb8tun
 //albuto

pipeline {
    agent any

    environment {
        GO_VERSION = '1.21'
        GO111MODULE = 'on'
        GOPROXY = 'https://proxy.golang.org'
        APP_NAME = 'jalitalks'
    }

    tools {
        go "${GO_VERSION}" // Jenkins global tools me Go 1.21 configure hona chahiye
    }

    stages {
        stage('Checkout') {
            steps {
                echo '🔄 Checking out source code...'
                checkout scm
            }
        }

        stage('Install Dependencies') {
            steps {
                echo '📦 Installing Go module dependencies...'
                sh 'go mod tidy'
            }
        }

        stage('Build') {
            steps {
                echo '🏗️ Building the Go application...'
                sh 'go build -o bin/${APP_NAME} main.go'
            }
        }

        stage('Test') {
            steps {
                echo '🧪 Running unit tests...'
                sh 'go test ./... -v -cover'
            }
        }

        stage('Docker Build') {
            steps {
                echo '🐳 Building Docker image...'
                script {
                    def imageTag = "${APP_NAME}:latest"
                    sh "docker build -t ${imageTag} ."
                }
            }
        }

        stage('Archive Build Artifact') {
            steps {
                echo '📁 Archiving binary for Jenkins...'
                archiveArtifacts artifacts: 'bin/**', fingerprint: true
            }
        }
    }

    post {
        always {
            echo '🏁 Pipeline finished!'
        }
        failure {
            echo '❌ Build failed. Check the logs.'
        }
        success {
            echo '✅ Build & Test successful!'
        }
    }
}
