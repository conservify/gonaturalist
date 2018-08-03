timestamps {
    node () {
        stage ('git') {
            checkout([$class: 'GitSCM', branches: [[name: '*/master']], userRemoteConfigs: [[url: 'https://github.com/Conservify/gonaturalist.git']]])
        }

        stage ('build') {
            sh "make clean deps all"
	      }
    }
}
