# orchestrate

Project for learning go by building an orchestrator (like k8s).
Let's see how far I get.

The components of an orchestrator:
1. task
2. the job
3. scheduler
4. manager
5. worker
6. cluster
7. storage
8. command line interface (cli)

Different systems will have different names for these components.
- Kubernetes calls the manager 'control plane' and worker 'kubelet'.
- Borg called the manager 'borgmaster' and worker 'borglet'.
- Hashicorp's Nomad calls the manager 'server', the worker 'client'.
etc.

Each of the manager's and worker have a storage.
The manager talks to the scheduler and the worker, but the workers and scheduler never talk except through the manager.

The term container is shorthand for process and resource isolation.
It is namespaces and control groups (cgroups) which are features of linux kernel.
Namespaces are mechanisms to isolate processes and their resources from each other.
Cgroups provide limits and accounting for a collection of processes.

With containers, an application can be decoupled from the os. 
So you can host multiple applications in their own container, 
and each application listens on the same port. 
In fact, you could have the same application running across multiple 
containers all using the same port. You would not have to give each of them 
different ports. The benefit of containers is that they give the application
the impression that they are the only application running on the os, and thus have all the
resources to themselves.

You can use an orchestrator to deploy and manage the application. The orchestrator automates
deployment, scaling and managing of the containers. Orchestrators are like CPU schedulers, but 
instead of managing os processes, it manages containers.


Tasks
-----
Tasks are the smallest unit of work and typically run in a container.
Tasks specify:
- cpu, memory and disk needed to run effectively
- restart policy - what should the orchestrator do in the event of failure
- name of the container image the task runs on

Job
---
Job is an aggregation of tasks. It has 1 or more tasks that form a larger 
grouping of taks to perform a set of function. For example a job maybe 
a RESTful API server and a reverse proxy.

Scheduler
---------
The scheduler decides what machine to run the tasks of a job.
This is simply selecting the node from a set of machines in a round robin, or 
it can be more complex selection process (based on a bunch of selection criteria).

Manager
------
The manager is the brain of the orchestrator and the main entry point
for the users. To run jobs, the humans submit their jobs to the manager.
The manage will use the scheduler to find a machine where the job's task can
run.  

The manager will also collect metrics, keep state of tasks, track the
managers tasks run on.

Workers
-------
Provider the muscle to the orchestrator. Responsible for running the tasks assigned by the manager.
If a task fails, the worker will attempt to restart a task.


Cluster
-------
The cluster is a logical grouping of all the previous components. 
Most of the time, clusters are run on separate physical machines, but they
can be on the same machine.

Clusters are used for high availability, scalability.
