timestamps {
    node () {
        dir ("../src/github.com/Conservify/gonaturalist") {
            stage ('git') {
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], userRemoteConfigs: [[url: 'https://github.com/Conservify/gonaturalist.git']]])
            }

            stage ('build') {
                sh "make clean deps all"
            }
        }
    }
}
