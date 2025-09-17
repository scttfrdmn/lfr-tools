# Educational Workflows: Real Classroom Examples

This guide shows how to use LFR Tools for different types of classes and educational scenarios.

## Computer Science Programming Class

### Class Setup: "Introduction to Python Programming"
**Students**: 30 beginners
**Duration**: 16-week semester
**Budget**: $1000 total

#### Week 1: Setup
```bash
# Teacher creates class environment
lfr students setup environment \
  --project=intro-python-fall2024 \
  --s3-bucket=intro-python-status \
  --start-date=2024-08-20 \
  --end-date=2024-12-10

# Create student computers (small size for beginners)
lfr users template students.csv
# Edit CSV with 30 student names, app_standard_xl_1_0 bundle

lfr users create-bulk students.csv --start-stopped

# Install Python environment on all computers
for student in alice bob charlie; do
  lfr software install python-dev $student
done

# Generate access tokens
lfr students generate tokens --project=intro-python-fall2024
```

#### Weekly Class Routine
```bash
# Before class (5 minutes early):
lfr instances start --project=intro-python-fall2024 --wait

# During class - check who needs help:
lfr students status --project=intro-python-fall2024

# After class - save money:
lfr instances stop --project=intro-python-fall2024
```

#### Student Experience
```bash
# Student connects to do homework:
lfr connect alice

# Instance auto-starts, student writes code:
alice@instance:~$ nano homework1.py
alice@instance:~$ python3 homework1.py
alice@instance:~$ exit

# Computer goes to sleep automatically after 2 hours of no activity
```

**Expected Costs**: ~$600 for full semester (under budget!)

---

## Data Science Research Lab

### Lab Setup: "Machine Learning Research"
**Students**: 8 graduate students
**Duration**: Full academic year
**Budget**: $3000 total, $300 per student

#### Initial Setup
```bash
# Create research environment
lfr students setup environment \
  --project=ml-research-lab \
  --s3-bucket=ml-lab-status \
  --students=researcher1,researcher2,researcher3 \
  --start-date=2024-09-01 \
  --end-date=2025-05-31

# Create powerful computers for data science
lfr users create \
  --project=ml-research-lab \
  --blueprint=ubuntu_22_04 \
  --bundle=app_standard_2xl_1_0 \
  --region=us-west-2 \
  --users=researcher1,researcher2,researcher3

# Install data science environment
for researcher in researcher1 researcher2 researcher3; do
  lfr software install data-science $researcher
done

# Create shared storage for datasets
lfr efs create lab-datasets --project=ml-research-lab
lfr efs mount-all fs-dataset123 --project=ml-research-lab
```

#### Daily Usage
```bash
# Researchers work independently:
lfr connect researcher1

# Long-running jobs (8+ hours) use research idle policy:
lfr idle advanced policies apply research-long-running --users=researcher1

# Shared data analysis:
researcher1@instance:~$ ls /mnt/efs/shared/datasets/
researcher1@instance:~$ python3 analyze_data.py /mnt/efs/shared/datasets/experiment1.csv
```

**Expected Costs**: ~$150/month per active researcher

---

## Web Development Bootcamp

### Class Setup: "Full Stack Web Development"
**Students**: 50 students
**Duration**: 12-week intensive program
**Budget**: $2000 total

#### Setup for High-Intensity Learning
```bash
# Create bootcamp environment
lfr students setup environment \
  --project=webdev-bootcamp-2024 \
  --s3-bucket=webdev-status \
  --students=@bootcamp-students.csv

# Create medium-sized computers for web development
lfr users create-bulk students.csv --start-stopped

# Install web development tools
for student in $(cat student-list.txt); do
  lfr software install web-dev $student
done

# Use aggressive cost optimization (students work in focused sprints)
lfr idle advanced policies apply development-aggressive --project=webdev-bootcamp-2024
```

#### Daily Bootcamp Schedule
```bash
# 9:00 AM - Prep for class
lfr instances start --project=webdev-bootcamp-2024

# 12:00 PM - Lunch break (optional cost saving)
lfr instances stop --project=webdev-bootcamp-2024

# 1:00 PM - Resume afternoon session
lfr instances start --project=webdev-bootcamp-2024

# 6:00 PM - End of day
lfr instances stop --project=webdev-bootcamp-2024
```

**Expected Costs**: ~$1200 for 12 weeks (under budget with aggressive optimization)

---

## Research Computing Workshop

### Workshop Setup: "Introduction to Cloud Computing"
**Participants**: 15 researchers from different departments
**Duration**: 2-day workshop
**Budget**: $200 total

#### Quick Workshop Setup
```bash
# Create temporary workshop environment
lfr students setup environment \
  --project=cloud-workshop-nov2024 \
  --s3-bucket=workshop-status \
  --start-date=2024-11-15 \
  --end-date=2024-11-17

# Create small computers for demonstration
lfr users create \
  --project=cloud-workshop-nov2024 \
  --blueprint=ubuntu_22_04 \
  --bundle=app_standard_xl_1_0 \
  --region=us-west-2 \
  --users=workshop1,workshop2,workshop3

# Pre-install common tools
for user in workshop1 workshop2 workshop3; do
  lfr software install python-dev $user
done
```

#### Workshop Day Management
```bash
# Day 1: 9:00 AM - Start workshop
lfr instances start --project=cloud-workshop-nov2024 --wait

# Participants connect during workshop:
# Each person gets: lfr connect workshop1, workshop2, etc.

# Day 2: 5:00 PM - Workshop ends
lfr instances stop --project=cloud-workshop-nov2024

# Cleanup after workshop:
lfr users remove-bulk --project=cloud-workshop-nov2024 --confirm
```

**Expected Costs**: ~$50 for 2 days

---

## Advanced Graduate Seminar

### Seminar Setup: "High Performance Computing"
**Students**: 6 PhD students
**Duration**: 16 weeks
**Budget**: $5000 (high-performance needs)

#### Setup for Intensive Computing
```bash
# Create high-performance environment
lfr students setup environment \
  --project=hpc-seminar-2024 \
  --s3-bucket=hpc-status

# Create powerful GPU computers
lfr users create \
  --project=hpc-seminar-2024 \
  --blueprint=ubuntu_22_04 \
  --bundle=gpu_nvidia_2xl_1_0 \
  --region=us-west-2 \
  --users=phd1,phd2,phd3,phd4,phd5,phd6

# Install GPU machine learning environment
for student in phd1 phd2 phd3 phd4 phd5 phd6; do
  lfr software install gpu-ml $student
done

# Create large shared storage for datasets
lfr efs create hpc-datasets --project=hpc-seminar-2024

# Use research-friendly idle detection (8-hour threshold)
lfr idle advanced policies apply research-long-running --project=hpc-seminar-2024
```

#### Research Usage Pattern
```bash
# Students work on long-running computations
phd1@instance:~$ python3 train_neural_network.py  # Runs for 6 hours
phd2@instance:~$ ./molecular_simulation.sh        # Runs overnight

# Shared data access:
phd3@instance:~$ cp /mnt/efs/shared/large_dataset.zip ~/
phd3@instance:~$ unzip large_dataset.zip
phd3@instance:~$ python3 analyze_data.py
```

**Expected Costs**: ~$300/month per student for GPU usage

---

## Mixed-Level Computer Science Course

### Course Setup: "Data Structures and Algorithms"
**Students**: 40 students (mixed experience levels)
**Duration**: 16 weeks
**Budget**: $1500 total

#### Tiered Resource Allocation
```bash
# Create different computer sizes based on student needs
# Beginners get smaller computers, advanced students get more powerful ones

# Create CSV with mixed bundles:
# alice,cs201,ubuntu_22_04,app_standard_xl_1_0,beginners
# bob,cs201,ubuntu_22_04,app_standard_2xl_1_0,advanced
# charlie,cs201,ubuntu_22_04,app_standard_xl_1_0,beginners

lfr users create-bulk mixed-students.csv --start-stopped

# Install appropriate software for different needs
lfr software install python-dev alice     # Beginner gets basic Python
lfr software install data-science bob     # Advanced gets full data science
```

#### Flexible Class Management
```bash
# Different idle policies for different groups:
lfr idle advanced policies apply educational-conservative --users=alice,charlie
lfr idle advanced policies apply educational-balanced --users=bob

# Monitor and adjust:
lfr idle advanced analyze --project=cs201
# Shows which students need more/less powerful computers
```

---

## Summer Research Program

### Program Setup: "Undergraduate Research Experience"
**Students**: 20 undergraduates
**Duration**: 10 weeks
**Budget**: $2000

#### Research-Focused Environment
```bash
# Create summer program environment
lfr students setup environment \
  --project=summer-research-2024 \
  --s3-bucket=summer-research-status \
  --start-date=2024-06-01 \
  --end-date=2024-08-10

# Create research-sized computers
lfr users create-bulk research-students.csv \
  --start-stopped \
  --from-snapshot=research-template-snapshot

# Set up collaborative environment
lfr efs create research-shared --project=summer-research-2024
lfr efs mount-all fs-research123 --project=summer-research-2024 --mode=rw

# Use balanced idle detection for research work
lfr idle advanced policies apply educational-balanced --project=summer-research-2024
```

#### Research Workflow
```bash
# Students work on independent projects but share resources:
student1@instance:~$ git clone https://github.com/lab/research-project
student1@instance:~$ cd research-project
student1@instance:~$ python3 experiment.py --data=/mnt/efs/shared/datasets/

# Collaborative analysis:
student1@instance:~$ cp results.csv /mnt/efs/shared/team-results/
student2@instance:~$ python3 combine_results.py /mnt/efs/shared/team-results/*.csv
```

---

## Key Patterns for Success

### Start Small and Scale
1. **Test with 3-5 students first**
2. **Work out the problems**
3. **Then scale to full class size**

### Cost Management
1. **Always use `--start-stopped` for new computers**
2. **Stop computers after class**
3. **Monitor usage weekly**
4. **Apply appropriate idle detection policies**

### Student Support
1. **Provide clear, simple instructions**
2. **Test everything from a student's perspective**
3. **Have TAs or advanced students help others**
4. **Keep backup plans for technical failures**

### File Organization
1. **Use shared storage for class materials**
2. **Create clear folder structures**
3. **Set appropriate permissions (read-only vs read-write)**
4. **Back up important work regularly**

These examples show how LFR Tools adapts to different educational needs while keeping costs reasonable and management simple. Start with a pattern similar to your class and customize as needed! 🎓