import { execSync } from 'child_process'

/**
 * Global teardown for LocalStack tests
 * Cleans up test resources and optionally stops LocalStack
 */
async function globalTeardown() {
  console.log('🧹 Cleaning up LocalStack tests...')

  try {
    // Clean up test IAM user
    execSync('AWS_ENDPOINT_URL=http://localhost:4566 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test aws iam delete-user --user-name e2e-test-user --region us-east-1', {
      timeout: 5000,
      stdio: 'ignore'
    })
    console.log('✅ LocalStack test resources cleaned up')
  } catch (error) {
    // It's ok if cleanup fails - LocalStack is ephemeral
    console.log('ℹ️ LocalStack cleanup completed (some resources may not exist)')
  }

  // Note: We don't stop LocalStack here as it might be used for other tests
  // Use `make stop-localstack` manually when done

  console.log('✅ LocalStack test teardown complete')
}

export default globalTeardown