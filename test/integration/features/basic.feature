@basic
Feature: Basic
  In order to use Minishift
  As a user
  I need to be able to bring up a test environment

  Scenario: Starting Minishift
    Given Minishift has state "Does Not Exist"
     When executing "minishift start --docker-env=FOO=BAR --docker-env=BAZ=BAT"
     Then Minishift should have state "Running"
      And Minishift should have a valid IP address

  Scenario: OpenShift developer account has sudo permissions
     The 'developer' user should be configured with the sudoer role after starting Minishift
     When executing "oc --as system:admin get clusterrolebindings"
     Then stderr should be empty
      And exitcode should equal 0
      And stdout should contain
     """
     sudoer
     """

  Scenario: A 'minishift' context is created for 'oc' usage
    After a successful Minishift start the user's current context is 'minishift'
    When executing "oc config current-context"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
    """
    minishift
    """

  Scenario: User can switch the current 'oc' context and return to 'minishift' context
    Given executing "oc config set-context dummy"
      And executing "oc config use-context dummy"
     When executing "oc project -q"
     Then exitcode should equal 1
     When executing "oc config use-context minishift"
      And executing "oc config current-context"
     Then stderr should be empty
      And exitcode should equal 0
      And stdout should contain
    """
    minishift
    """

  Scenario: User has a pre-configured set of persitence volumnes
    When executing "oc get pv --as system:admin -o=name"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     persistentvolume/pv0001
     """

  Scenario: User is able to do ssh into Minishift VM
    Given Minishift has state "Running"
     When executing "minishift ssh echo hello"
     Then stderr should be empty
      And exitcode should equal 0
      And stdout should contain
      """
      hello
      """

  Scenario: User is able to set custom Docker specific environment variables
    Given Minishift has state "Running"
     When executing "minishift ssh cat /var/lib/boot2docker/profile"
     Then stderr should be empty
      And exitcode should equal 0
      And stdout should contain
      """
      export "FOO=BAR"
      export "BAZ=BAT"
      """

  Scenario: User is able to retrieve host and port of OpenShift registry
    Given Minishift has state "Running"
     When executing "minishift openshift registry"
     Then stderr should be empty
      And exitcode should equal 0
      And stdout should contain
      """
      172.30.1.1:5000
      """

  # User can deploy the example Ruby application ruby-ex
  Scenario: User can login to the server
    When executing "oc login --username=developer --password=developer"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     Login successful
     """
  
  Scenario: User can create new namespace ruby for application ruby-ex
    When executing "oc new-project ruby"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     Now using project "ruby"
     """
  
  Scenario: User can deploy application ruby-ex to namespace ruby 
    When executing "oc new-app centos/ruby-22-centos7~https://github.com/openshift/ruby-ex.git"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     Success
     """
    When executing "oc rollout status deploymentconfig ruby-ex --watch"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     "ruby-ex-1" successfully rolled out
     """
  
  Scenario: User can create route for ruby-ex to make it visiable outside of the cluster
    When executing "oc expose svc/ruby-ex"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     exposed
     """
  
  Scenario: User can delete namespace ruby
    When executing "oc delete project ruby"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     "ruby" deleted
     """
  
  Scenario: User can log out the session
    When executing "oc logout"
    Then stderr should be empty
     And exitcode should equal 0
     And stdout should contain
     """
     Logged "developer" out
     """
  # End of Ruby application ruby-ex deployment
  
  Scenario: Stopping Minishift
    Given Minishift has state "Running"
     When executing "minishift stop"
     Then Minishift should have state "Stopped"

  Scenario: Deleting Minishift
    Given Minishift has state "Stopped"
     When executing "minishift delete"
     Then Minishift should have state "Does Not Exist"
