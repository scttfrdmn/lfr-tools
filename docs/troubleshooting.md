# Troubleshooting Guide

This guide helps you solve common problems with LFR Tools using simple language.

## For Students

### Connection Problems

#### "No access token found"

**What this means**: Your computer doesn't know how to connect to your cloud computer.

**How to fix it:**
1. Make sure you ran the setup command your teacher gave you
2. Check that you typed your username correctly
3. Try the activate command again:
   ```bash
   lfr connect activate <your-token> <your-student-id>
   ```

#### "Instance is not running"

**What this means**: Your cloud computer is turned off (sleeping).

**How to fix it:**
1. Wait 1-2 minutes - it might be starting up
2. Try connecting again: `lfr connect alice`
3. If it still doesn't work, ask your teacher to wake up the computers

#### "Connection failed" or "Connection timed out"

**What this means**: There's a problem connecting to your cloud computer.

**How to fix it:**
1. Check your internet connection
2. Make sure you typed your username correctly
3. Wait a few minutes and try again
4. Ask your teacher if there are any known problems

#### "Token is bound to a different machine"

**What this means**: You're trying to use your access token on a different computer.

**How to fix it:**
1. Use the same computer where you first set up access
2. If you got a new computer, ask your teacher for a new token
3. If you need to use a different computer temporarily, ask your teacher for help

### Using Your Cloud Computer

#### "Permission denied" when accessing files

**What this means**: You're trying to access files you don't have permission to read or change.

**How to fix it:**
1. Make sure you're in your own folder: `cd ~`
2. Don't try to access other students' folders
3. For shared files, use: `/mnt/efs/shared/`
4. Ask your teacher about file permissions

#### "Command not found"

**What this means**: The program you're trying to run isn't installed.

**How to fix it:**
1. Check the spelling of the command
2. Ask your teacher if the software is installed
3. Try using the full path: `/usr/bin/python3` instead of `python3`

#### "No space left on device"

**What this means**: Your cloud computer's storage is full.

**How to fix it:**
1. Delete old files you don't need: `rm old_file.txt`
2. Check what's taking up space: `ls -la`
3. Ask your teacher about getting more storage
4. Move large files to shared storage

### File and Programming Issues

#### "File not found"

**What this means**: The file you're looking for doesn't exist or isn't in the current folder.

**How to fix it:**
1. Check what files are in the current folder: `ls`
2. Make sure you spelled the filename correctly
3. Check if you're in the right folder: `pwd`
4. Look for the file: `find ~ -name "myfile.py"`

#### "Syntax error" in your code

**What this means**: There's a mistake in your program code.

**How to fix it:**
1. Read the error message carefully - it usually tells you the line number
2. Check for missing punctuation: `:`, `;`, `,`, `)`
3. Make sure your indentation is correct (especially in Python)
4. Ask your teacher or classmates for help

## For Teachers

### Setup Problems

#### "Failed to create AWS client"

**What this means**: LFR Tools can't connect to AWS.

**How to fix it:**
1. Check your AWS credentials: `aws sts get-caller-identity`
2. Make sure you're using the right profile: `--profile=aws`
3. Check your internet connection
4. Contact your IT department if AWS access is blocked

#### "S3 bucket creation failed"

**What this means**: Can't create the communication system for students.

**How to fix it:**
1. Choose a different bucket name (must be globally unique)
2. Check AWS permissions for S3 access
3. Use an existing bucket if you have one
4. Contact AWS support if permissions are blocked

#### "User creation failed"

**What this means**: Can't create student accounts.

**How to fix it:**
1. Check AWS IAM permissions
2. Make sure you have enough AWS service limits
3. Try creating one user at a time to identify the problem
4. Check the error message for specific details

### Student Access Problems

#### "Students can't connect"

**What this means**: The access system isn't working properly.

**How to fix it:**
1. Check that you ran: `lfr students setup environment`
2. Make sure you generated tokens: `lfr students generate tokens`
3. Verify the S3 bucket exists and has the right permissions
4. Test with your own test account first

#### "Too many start requests"

**What this means**: Many students are trying to wake up their computers at once.

**How to fix it:**
1. Start all computers before class: `lfr instances start --project=cs101 --wait`
2. Use auto-approve during class hours: `lfr students check requests --auto-approve`
3. Consider pre-warming computers before class time

#### "Costs are too high"

**What this means**: Your AWS bill is higher than expected.

**How to fix it:**
1. Check usage patterns: `lfr idle advanced analyze --project=cs101`
2. Apply cost-saving policies: `lfr idle advanced policies apply educational-conservative`
3. Stop computers after class: `lfr instances stop --project=cs101`
4. Consider smaller computer sizes for basic tasks

### Class Management Issues

#### "Student lost their work"

**What this means**: Files are missing from a student's computer.

**How to fix it:**
1. Connect to their computer: `lfr ssh connect alice`
2. Look in their home folder: `ls /home/alice/`
3. Check shared folders: `ls /mnt/efs/shared/`
4. Check if they saved files in a different location

#### "Software installation failed"

**What this means**: Couldn't install required software on student computers.

**How to fix it:**
1. Make sure the computer is running: `lfr instances start --users=alice`
2. Check if the software pack is compatible: `lfr software install python-dev alice --force`
3. Try installing manually via SSH: `lfr ssh connect alice`
4. Create a custom software pack if needed

#### "Student can't access shared files"

**What this means**: File permissions or mounting problems.

**How to fix it:**
1. Check if EFS is set up: `lfr efs status`
2. Make sure EFS is mounted: `lfr efs mount-status --project=cs101`
3. Connect to student's computer and check: `ls /mnt/efs/`
4. Recreate the EFS mount if needed

## Common Error Messages

### "Operation timed out"
- **Problem**: AWS is taking too long to respond
- **Solution**: Wait a few minutes and try again
- **Prevention**: Use `--wait` flags for long operations

### "Insufficient permissions"
- **Problem**: Your AWS account can't perform the action
- **Solution**: Check with your IT department about AWS permissions
- **Prevention**: Test with a small setup first

### "Budget exhausted"
- **Problem**: You've reached your spending limit
- **Solution**: Check costs and extend budget if needed
- **Prevention**: Set up cost monitoring and alerts

### "Instance not found"
- **Problem**: Looking for a computer that doesn't exist
- **Solution**: Check the username and project name spelling
- **Prevention**: Use `lfr instances list` to see what exists

## Best Practices

### For Reliable Classes
1. **Test everything** before the first class
2. **Start computers early** - don't wait until class begins
3. **Have backup plans** - technology sometimes fails
4. **Monitor costs** - check spending weekly
5. **Keep it simple** - start with basic features, add complexity slowly

### For Student Success
1. **Provide clear instructions** - students need step-by-step guides
2. **Test student workflows** - try everything from a student's perspective
3. **Have helpers available** - TAs or advanced students can help others
4. **Practice problem-solving** - know how to fix common issues quickly

### For Cost Control
1. **Use `--start-stopped`** when creating computers
2. **Stop computers after class** - don't leave them running overnight
3. **Monitor usage patterns** - identify computers that are always idle
4. **Right-size computers** - don't use GPU computers for basic tasks

## Emergency Procedures

### "Class starts in 10 minutes and nothing works"

**Immediate actions:**
1. **Stay calm** - most problems have quick fixes
2. **Check basics**: Internet connection, AWS access
3. **Start backup plan**: Use local computers or postpone technical work
4. **Get help**: Contact IT department or someone who knows AWS

**Quick diagnostics:**
```bash
# Check if LFR Tools is working:
lfr --version

# Check AWS connection:
lfr instances list

# Check student status:
lfr students status --project=cs101
```

### "Student deadlines approaching and computers aren't working"

**Immediate actions:**
1. **Extend deadlines** if needed - technology problems happen
2. **Start all computers**: `lfr instances start --project=cs101 --wait`
3. **Check individual student issues**: `lfr students status --project=cs101`
4. **Provide alternative access** if needed

### "Unexpected high costs"

**Immediate actions:**
1. **Check what's running**: `lfr instances list --project=cs101`
2. **Stop unnecessary computers**: `lfr instances stop --project=cs101`
3. **Check usage patterns**: `lfr idle advanced analyze --project=cs101`
4. **Apply cost controls**: `lfr idle advanced policies apply educational-conservative`

## Getting Additional Help

### Self-Help Resources
- **Built-in help**: Add `--help` to any command
- **Debug mode**: Add `--debug` to see detailed information
- **Documentation**: Check the docs folder for detailed guides

### Community Support
- **GitHub Issues**: Report bugs and get help
- **Documentation**: Complete guides for all features
- **Examples**: Sample configurations for common scenarios

### Professional Support
- **AWS Support**: For underlying infrastructure issues
- **IT Department**: For institutional AWS account problems
- **LFR Tools Community**: For feature requests and usage questions

Remember: Every teacher starts as a beginner with cloud technology. Don't be afraid to experiment with test accounts first! 🎓