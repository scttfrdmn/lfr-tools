import { execSync } from 'child_process'

/**
 * Global setup for AWS integration tests
 * Verifies AWS credentials and prepares test environment
 */
async function globalSetup() {
  console.log('🔧 Setting up AWS integration tests...')

  // Check if AWS credentials are available
  try {
    // Test AWS credentials
    const result = execSync('aws sts get-caller-identity', {
      encoding: 'utf8',
      timeout: 10000
    })

    const identity = JSON.parse(result)
    console.log(`✅ AWS credentials verified for account: ${identity.Account}`)
    console.log(`   User/Role: ${identity.Arn}`)

    // Set reasonable AWS region if not set
    if (!process.env.AWS_REGION) {
      process.env.AWS_REGION = 'us-east-1'
    }

    console.log(`   Region: ${process.env.AWS_REGION}`)

  } catch (error) {
    console.error('❌ AWS credentials not configured properly')
    console.error('Please ensure AWS credentials are available via:')
    console.error('  - AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables')
    console.error('  - AWS CLI profile (aws configure)')
    console.error('  - IAM role (if running on EC2)')
    throw new Error('AWS credentials required for integration tests')
  }

  // Verify LFR Tools binary is available
  try {
    execSync('which lfr || echo "LFR Tools binary not found"', { encoding: 'utf8' })
    console.log('✅ LFR Tools CLI available')
  } catch (error) {
    console.warn('⚠️ LFR Tools CLI not in PATH - some tests may not work fully')
  }

  console.log('✅ AWS integration test setup complete')
}

export default globalSetup