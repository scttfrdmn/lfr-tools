# Common Tasks: Step-by-Step Examples

This guide shows you how to do common tasks with simple, real examples.

## For Students

### Task 1: Connect to Your Cloud Computer

**Goal**: Get access to your assigned computer for class work.

**Steps:**
1. Open your terminal (Terminal on Mac, Command Prompt on Windows)
2. Type your connection command:
   ```bash
   lfr connect alice
   ```
   (Replace "alice" with your actual username)

3. Wait for the connection:
   ```
   Connecting to alice's instance in project cs101...
   Instance state: stopped
   Requesting start from instructor...
   ⏳ Waiting for instance to start (elapsed: 30s)
   ✅ Connected!
   ```

4. You now see a new prompt like: `alice@instance:~$`
5. You're now using your cloud computer!

### Task 2: Save and Run a Python Program

**Goal**: Write a simple Python program and run it.

**Steps:**
1. Connect to your cloud computer (see Task 1)
2. Create a new Python file:
   ```bash
   nano hello.py
   ```
3. Type this simple program:
   ```python
   print("Hello from the cloud!")
   print("My name is Alice")
   ```
4. Save the file:
   - Press `Ctrl + X`
   - Press `Y` (yes, save)
   - Press `Enter` (keep the filename)

5. Run your program:
   ```bash
   python3 hello.py
   ```
6. You should see:
   ```
   Hello from the cloud!
   My name is Alice
   ```

### Task 3: Access Shared Class Files

**Goal**: Get files your teacher shared with the class.

**Steps:**
1. Connect to your cloud computer
2. Look at shared files:
   ```bash
   ls /mnt/efs/shared/
   ```
3. You might see:
   ```
   assignment1.py  lecture_notes.txt  sample_data.csv
   ```
4. Copy a file to your folder:
   ```bash
   cp /mnt/efs/shared/assignment1.py ~/
   ```
5. Now you can edit it:
   ```bash
   nano assignment1.py
   ```

### Task 4: Submit Your Homework

**Goal**: Turn in your completed assignment.

**Steps:**
1. Make sure your work is saved in your home folder
2. Copy your file to the submissions folder:
   ```bash
   cp my_assignment.py /mnt/efs/submissions/alice_assignment1.py
   ```
   (Include your name in the filename)
3. Verify it was submitted:
   ```bash
   ls /mnt/efs/submissions/alice*
   ```
4. You should see your file listed

### Task 5: Install Additional Software

**Goal**: Add software you need for your project.

**Steps:**
1. Connect to your cloud computer
2. See what's available to install:
   ```bash
   apt list | grep python
   ```
3. Install software (example - installing a Python package):
   ```bash
   pip3 install --user requests
   ```
4. Test that it works:
   ```bash
   python3 -c "import requests; print('Requests package installed!')"
   ```

## For Teachers

### Task 1: Set Up a New Class

**Goal**: Create cloud computers for all students in your class.

**Steps:**
1. Create a list of students in a file called `students.csv`:
   ```csv
   username,project,blueprint,bundle,groups
   alice,cs101,ubuntu_22_04,app_standard_xl_1_0,students
   bob,cs101,ubuntu_22_04,app_standard_xl_1_0,students
   charlie,cs101,ubuntu_22_04,app_standard_xl_1_0,students
   ```

2. Set up the class environment:
   ```bash
   lfr students setup environment \
     --project=cs101-fall2024 \
     --s3-bucket=cs101-status-bucket \
     --students=alice,bob,charlie
   ```

3. Create the student computers:
   ```bash
   lfr users create-bulk students.csv --start-stopped
   ```

4. Generate access codes for students:
   ```bash
   lfr students generate tokens --project=cs101-fall2024
   ```

5. Send each student their access code via email

### Task 2: Start Computers for Class

**Goal**: Wake up all student computers before class begins.

**Steps:**
1. Check current status:
   ```bash
   lfr instances list --project=cs101-fall2024
   ```
2. Start all computers:
   ```bash
   lfr instances start --project=cs101-fall2024 --wait
   ```
3. You'll see:
   ```
   Starting 25 instances for users: [alice,bob,charlie...]
   ⏳ Waiting for instances to start...
   ✅ All instances started after 2m30s
   ```
4. Tell students they can now connect

### Task 3: Install Software for Everyone

**Goal**: Set up Python programming environment for all students.

**Steps:**
1. See available software packs:
   ```bash
   lfr software list
   ```
2. Install Python development tools on one student's computer (test first):
   ```bash
   lfr software install python-dev alice
   ```
3. If it works, install for everyone:
   ```bash
   # You'll need to run this for each student
   lfr software install python-dev bob
   lfr software install python-dev charlie
   # etc.
   ```

### Task 4: Help a Struggling Student

**Goal**: Connect to a student's computer to help them debug their code.

**Steps:**
1. Make sure their computer is running:
   ```bash
   lfr instances list --user=alice
   ```
2. If stopped, start it:
   ```bash
   lfr instances start --users=alice --wait
   ```
3. Connect to their computer:
   ```bash
   lfr ssh connect alice
   ```
4. You're now on their computer and can help debug
5. When done, type `exit` to return to your computer

### Task 5: Set Up Shared Class Files

**Goal**: Create a shared folder where you can put class materials and students can submit work.

**Steps:**
1. Create shared storage:
   ```bash
   lfr efs create class-shared --project=cs101-fall2024
   ```
2. Note the filesystem ID (like `fs-1234567890abcdef0`)
3. Mount it on all student computers:
   ```bash
   lfr efs mount-all fs-1234567890abcdef0 --project=cs101-fall2024
   ```
4. Connect to a computer and set up folders:
   ```bash
   lfr ssh connect alice
   sudo mkdir -p /mnt/efs/shared/assignments
   sudo mkdir -p /mnt/efs/shared/examples
   sudo mkdir -p /mnt/efs/submissions
   sudo chmod 755 /mnt/efs/shared/
   sudo chmod 733 /mnt/efs/submissions/  # Students can submit but not see others
   exit
   ```

### Task 6: Check Costs and Optimize

**Goal**: Make sure you're not spending too much money.

**Steps:**
1. Analyze current usage:
   ```bash
   lfr idle advanced analyze --project=cs101-fall2024
   ```
2. You'll see a report like:
   ```
   Current cost: $450/month
   With optimization: $225/month
   Total savings: $225/month (50% reduction)
   ```
3. Apply cost-saving policies:
   ```bash
   lfr idle advanced policies apply educational-conservative --project=cs101-fall2024
   ```
4. Stop computers when not needed:
   ```bash
   lfr instances stop --project=cs101-fall2024
   ```

## For IT Administrators

### Task 1: Set Up AWS Account for Education

**Goal**: Prepare AWS account for educational use.

**Steps:**
1. Create AWS account or use existing institutional account
2. Set up billing alerts to prevent surprise costs
3. Create IAM user for the teacher with appropriate permissions
4. Configure AWS CLI with teacher's credentials
5. Test basic LFR Tools functionality

### Task 2: Configure Budget Limits

**Goal**: Prevent runaway costs.

**Steps:**
1. Set up AWS Budget alerts in the AWS Console
2. Configure institutional spending limits
3. Set up automatic notifications for cost thresholds
4. Test budget enforcement with small test projects

### Task 3: Security Review

**Goal**: Ensure educational use meets security requirements.

**Steps:**
1. Review IAM permissions for teacher accounts
2. Ensure student tokens cannot access AWS directly
3. Verify S3 bucket policies are appropriately restricted
4. Test token sharing prevention (hardware binding)
5. Document access controls for compliance

## Quick Reference

### Most Common Student Commands
```bash
lfr connect alice          # Connect to your computer
lfr connect list          # See your available computers
exit                      # Disconnect from cloud computer
```

### Most Common Teacher Commands
```bash
lfr instances start --project=cs101 --wait    # Start class computers
lfr instances stop --project=cs101            # Stop class computers
lfr students status --project=cs101           # Check student status
lfr ssh connect alice                         # Help a specific student
```

### Emergency Commands
```bash
lfr instances list --project=cs101            # See all computers
lfr students check requests --auto-approve    # Approve student access
lfr instances start --users=alice --wait      # Start one computer
```

Remember: When in doubt, add `--help` to any command to see what it does!

## What to Do When You're Stuck

1. **Read the error message carefully** - it usually tells you what's wrong
2. **Try the command again** - sometimes it's a temporary problem
3. **Check your internet connection** - cloud tools need internet
4. **Ask for help** - teachers, TAs, classmates, or IT support
5. **Check the troubleshooting guide** - you might not be the first person with this problem

Don't give up! Cloud computing can seem complicated at first, but these tools are designed to make it as simple as possible. 🌟