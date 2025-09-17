# Real Examples: See LFR Tools in Action

This page shows real examples of using LFR Tools with actual commands and outputs.

## Example 1: Setting Up a Programming Class

**Scenario**: Ms. Johnson wants to set up cloud computers for her 25 students in "Intro to Programming."

### Teacher's Setup Process

**Step 1: Create the class environment**
```bash
$ lfr students setup environment \
  --project=intro-programming-2024 \
  --s3-bucket=intro-prog-status \
  --students=alice,bob,charlie,diana,emma \
  --professor=msjohnson \
  --start-date=2024-09-01 \
  --end-date=2024-12-15

Setting up class environment for project: intro-programming-2024
S3 bucket: intro-prog-status
Students: 5, TAs: 0, Professor: msjohnson
✅ S3 bucket created: intro-prog-status
✅ S3 sync enabled for project intro-programming-2024
✅ Class setup completed!

Next steps:
1. Generate tokens: lfr students generate tokens --project=intro-programming-2024
2. Create users: lfr users create-bulk students.csv
3. Distribute tokens to students
```

**Step 2: Create student computers**
```bash
$ lfr users template students.csv
✅ User CSV template created: students.csv

# Edit the CSV file to include all students, then:
$ lfr users create-bulk students.csv --start-stopped

Bulk user creation from: students.csv
Total users: 5

Creating users in 1 batch(es):

Batch 1: intro-programming-2024/ubuntu_22_04/app_standard_xl_1_0 (5 users)

[1/5] Creating user: alice
✅ alice : Kj8mN2pQ-Rt5vX9wZ : arn:aws:lightsail:us-west-2:123456789:Instance/alice-ubuntu_22_04

🛑 Stopping instances to save costs...
✅ Instances stopped for cost savings

🎉 Bulk user creation completed!
✅ Success: 5 users
```

**Step 3: Install Python on all computers**
```bash
$ lfr software install python-dev alice

Installing software pack: Python Development Environment
Description: Complete Python development setup with common packages
Target user: alice
Target instance: alice-ubuntu_22_04 (stopped)
Error: instance alice-ubuntu_22_04 is not running (state: stopped). Start it first

# Start Alice's computer first:
$ lfr instances start --users alice --wait
✅ instance alice-ubuntu_22_04 reached state 'running' after 45s

# Now install software:
$ lfr software install python-dev alice --force
✅ Software pack 'Python Development Environment' installed successfully!
📋 Installation script prepared for alice-ubuntu_22_04
```

### Student's First Connection

**Alice receives email from Ms. Johnson:**
```
Subject: Your Programming Class Cloud Computer

Hi Alice,

Welcome to Intro to Programming! Here's how to access your cloud computer:

1. Install LFR Tools: brew install lfr
2. Set up access: lfr connect activate intro-programming-2024-alice-Kj8mN2pQ 12345
3. Connect anytime: lfr connect alice

Your access is valid from Sept 1 to Dec 15.
See you in class!

Ms. Johnson
```

**Alice's setup:**
```bash
$ brew install lfr
🍺 lfr was successfully installed!

$ lfr connect activate intro-programming-2024-alice-Kj8mN2pQ 12345
Activating access token for student ID: 12345
Binding to current machine...
✅ Token activated!
Project: intro-programming-2024
Username: alice
Student ID: 12345

You can now connect with: lfr connect alice
```

**Alice's daily usage:**
```bash
$ lfr connect alice
Connecting to alice's instance in project intro-programming-2024...
Instance state: stopped
Requesting start from instructor...
⏳ ⠋ Waiting for instance to start (elapsed: 30s)
✅ Instance started and ready after 52s
Connecting to 34.221.15.92...

Welcome to Ubuntu 22.04.5 LTS
alice@instance:~$ python3 --version
Python 3.10.12

alice@instance:~$ nano hello.py
# Alice writes her first Python program

alice@instance:~$ python3 hello.py
Hello, World! My name is Alice.

alice@instance:~$ exit
Connection closed.
$
```

---

## Example 2: Research Lab with GPU Computing

**Scenario**: Dr. Smith's research lab needs GPU computers for machine learning research.

### Setting Up GPU Research Environment

```bash
$ lfr users create \
  --project=ai-research-lab \
  --blueprint=ubuntu_22_04 \
  --bundle=gpu_nvidia_2xl_1_0 \
  --region=us-west-2 \
  --users=grad_student1,grad_student2

Creating 2 users for project: ai-research-lab
Blueprint: ubuntu_22_04, Bundle: gpu_nvidia_2xl_1_0, Region: us-west-2

[1/2] Creating user: grad_student1
✅ grad_student1 : Pq7nM4tW-Yx3kL8vB : arn:aws:lightsail:us-west-2:123456789:Instance/grad_student1-ubuntu_22_04

[2/2] Creating user: grad_student2
✅ grad_student2 : Zf9rH5sC-Dt6wN2mX : arn:aws:lightsail:us-west-2:123456789:Instance/grad_student2-ubuntu_22_04

🎉 User creation completed!
```

### Installing GPU Software

```bash
$ lfr software install gpu-ml grad_student1
Installing software pack: GPU Machine Learning Environment
Description: CUDA, PyTorch, TensorFlow for GPU-enabled instances
Target user: grad_student1
Target instance: grad_student1-ubuntu_22_04 (running)
✅ Software pack 'GPU Machine Learning Environment' installed successfully!
```

### Research Usage

```bash
$ lfr connect grad_student1
grad_student1@instance:~$ nvidia-smi
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 535.86.10    Driver Version: 535.86.10    CUDA Version: 12.2   |
|-------------------------------+----------------------+----------------------+
|   0  Tesla T4           Off  | 00000000:00:1E.0 Off |                    0 |
+-----------------------------------------------------------------------------+

grad_student1@instance:~$ python3
>>> import torch
>>> torch.cuda.is_available()
True
>>> torch.cuda.get_device_name(0)
'Tesla T4'

grad_student1@instance:~$ python3 train_model.py
Training neural network on GPU...
Epoch 1/100: 95% accuracy
Epoch 2/100: 96% accuracy
...
```

---

## Example 3: Managing Costs and Budgets

**Scenario**: Professor needs to stay within a $500 semester budget for 30 students.

### Cost Analysis and Optimization

```bash
$ lfr idle advanced analyze --project=cs101-fall2024
Analyzing idle patterns for project: cs101-fall2024
Analysis period: 7 days

📊 Usage Pattern Analysis:
Instance          Avg CPU    Avg Memory    Network I/O    SSH Sessions    Recommended Policy
-----------------------------------------------------------------------------------------------
alice-ubuntu      2.1%       8.5%         Low           0.2/day        educational-conservative
bob-ubuntu        15.3%      25.1%        Medium        2.1/day        educational-balanced
charlie-gpu       45.2%      78.9%        High          4.5/day        research-long-running

💰 Cost Optimization Recommendations:
- Alice: Switch to educational-conservative (save ~40%)
- Bob: Apply educational-balanced (save ~60%)
- Charlie: Keep research-long-running (save ~25% safely)

📈 Projected Monthly Savings:
Current cost: $450/month
With optimization: $225/month
Total savings: $225/month (50% reduction)
```

### Applying Cost Controls

```bash
$ lfr idle advanced policies apply educational-conservative --users=alice
Applying idle detection policy: Educational Conservative
Category: educational
Estimated savings: 40%

✅ Policy applied to alice's instance
Idle detection: 3 hours threshold, 15 minute grace period
```

---

## Example 4: Troubleshooting Common Issues

### Student Can't Connect

**Student reports**: "I get 'no access token found' when I try to connect"

**Teacher troubleshooting:**
```bash
$ lfr students status --project=cs101-fall2024
Student status for project: cs101-fall2024

STUDENT         INSTANCE             STATE        PUBLIC IP          LAST ACTIVITY
------------------------------------------------------------------------------------------
alice           alice-ubuntu         running      34.221.15.92       Active now
bob             bob-ubuntu           stopped      -                  Stopped
charlie         charlie-ubuntu       stopped      -                  Stopped

# Alice's computer is running, so the problem is with her access token
```

**Solution**: Have Alice run the activate command again or provide a new token.

### Unexpected High Costs

**AWS bill shows $800 instead of expected $200**

**Teacher investigation:**
```bash
$ lfr instances list --project=cs101-fall2024
INSTANCE             STATE        PUBLIC IP          BLUEPRINT       BUNDLE               REGION       PROJECT
------------------------------------------------------------------------------------------------------------------------
alice-ubuntu         running      34.221.15.92       Ubuntu          app_standard_4xl_1_0  us-west-2    cs101-fall2024
bob-ubuntu           running      35.87.81.251       Ubuntu          gpu_nvidia_2xl_1_0   us-west-2    cs101-fall2024
charlie-ubuntu       running      52.12.98.142       Ubuntu          app_standard_4xl_1_0  us-west-2    cs101-fall2024

# Problem: Students are using expensive computer sizes!
```

**Solution**: Resize computers to appropriate sizes:
```bash
$ lfr instances resize alice-ubuntu down
$ lfr instances resize charlie-ubuntu down
$ lfr instances gpu bob-ubuntu disable  # Remove GPU if not needed
```

---

## Example 5: End of Semester Cleanup

**Goal**: Properly close a class and preserve student work.

```bash
# Week 16: Create final backups
$ lfr instances snapshot alice-ubuntu
$ lfr instances snapshot bob-ubuntu
$ lfr instances snapshot charlie-ubuntu

Creating snapshot: alice-ubuntu-snapshot
✅ Snapshot created: alice-ubuntu-snapshot

# Generate final cost report
$ lfr idle advanced analyze --project=cs101-fall2024 --days=120
📊 Semester Summary:
Total cost: $380 (under $500 budget!)
Average cost per student: $12.67
Most expensive student: Charlie ($25)
Least expensive student: Alice ($8)

# Archive shared files (optional)
$ lfr efs list --project=cs101-fall2024
EFS file systems for project: cs101-fall2024
NAME                 ID              STATE      MOUNT TARGETS
class-shared         fs-abc123       available  3

# Final cleanup (after grades are submitted)
$ lfr users remove-bulk --project=cs101-fall2024 --confirm
Bulk user removal:
Users to remove: [alice, bob, charlie, diana, emma]
Total: 5 users

⚠️  This will permanently delete:
   - 5 IAM users and their login profiles
   - All associated Lightsail instances
   - All user data and configurations

Removing 5 users...
✅ Success: 5 users

🎉 Bulk user removal completed!
```

These examples show how LFR Tools handles real classroom scenarios from setup to cleanup. Each example includes actual commands and expected outputs to help you understand what to expect! 📚