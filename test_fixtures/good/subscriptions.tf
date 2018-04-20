// DO NOT EDIT
// This file has been auto-generated by hedwig-terraform-generator {{GENERATOR_VERSION}}

module "sub-dev-myapp-my-topic" {
  source  = "Automatic/hedwig-queue-subscription/aws"
  version = "~> {{TFQueueSubscriptionModuleVersion}}"

  queue = "${module.consumer-dev-myapp.queue_arn}"
  topic = "${module.topic-my-topic.arn}"
}

module "sub-dev-myapp-my-topic2" {
  source  = "Automatic/hedwig-queue-subscription/aws"
  version = "~> {{TFQueueSubscriptionModuleVersion}}"

  queue = "${module.consumer-dev-myapp.queue_arn}"
  topic = "${module.topic-my-topic2.arn}"
}

module "sub-dev-secondapp-my-topic2" {
  source  = "Automatic/hedwig-queue-subscription/aws"
  version = "~> {{TFQueueSubscriptionModuleVersion}}"

  queue = "${module.consumer-dev-secondapp.queue_arn}"
  topic = "${module.topic-my-topic2.arn}"
}

module "sub-my-function-my-topic" {
  source  = "Automatic/hedwig-lambda-subscription/aws"
  version = "~> {{TFLambdaSubscriptionModuleVersion}}"

  function_arn       = "arn:aws:lambda:us-west-2:12345:function:myFunction:deployed"
  function_name      = "myFunction"
  function_qualifier = "deployed"
  topic              = "${module.topic-my-topic.arn}"
}