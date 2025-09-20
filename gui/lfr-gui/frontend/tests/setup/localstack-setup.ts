import { execSync } from 'child_process'

/**
 * Global setup for LocalStack tests
 * Ensures LocalStack is running and configured
 */
async function globalSetup() {
  console.log('🔧 Setting up LocalStack tests...')

  // Check if LocalStack is accessible
  try {
    // Wait for LocalStack to be ready
    let retries = 30
    let ready = false

    while (retries > 0 && !ready) {
      try {
        execSync('curl -s http://localhost:4566/_localstack/health', {
          timeout: 3000,
          stdio: 'ignore'
        })
        ready = true
      } catch {
        retries--
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
    }

    if (!ready) {
      throw new Error('LocalStack not accessible at http://localhost:4566')
    }

    console.log('✅ LocalStack is ready')

    // Initialize basic LocalStack resources for testing
    try {
      // Create a test IAM user to verify IAM operations work
      execSync('AWS_ENDPOINT_URL=http://localhost:4566 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws iam create-user --user-name e2e-test-user --region us-east-1', {
        encoding: 'utf8',
        timeout: 10000,
        stdio: 'ignore'
      })
      console.log('✅ LocalStack IAM test resources initialized')
    } catch (error) {
      // It's ok if the user already exists
      console.log('ℹ️ LocalStack resources may already exist')
    }

  } catch (error) {
    console.error('❌ LocalStack setup failed:', error)
    console.error('Please ensure LocalStack is running:')
    console.error('  make test-with-localstack')
    throw error
  }

  console.log('✅ LocalStack test setup complete')
}

export default globalSetup