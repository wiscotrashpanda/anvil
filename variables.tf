variable "stack_set_administration_role_arn" {
  description = "System-wide CloudFormation StackSet administration role ARN for self-managed StackSets."
  type        = string
  default     = null
}

variable "stack_set_execution_role_name" {
  description = "System-wide CloudFormation StackSet execution role name expected in target accounts."
  type        = string
  default     = "AWSCloudFormationStackSetExecutionRole"
}
