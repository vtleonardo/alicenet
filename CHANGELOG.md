# CHANGELOG

## **1.0.0-alpha-1 (unreleased)**

### **Features**

* Major refactoring in the layer 1 package (previously called blockchain package). Now the package is divided into 3 main components:

  * **Monitor**: responsible for monitoring and reacting tp layer 1 events (only ethereum at this point).
  * **Executor**: responsible for scheduling and executing tasks and actions (e.g snapshots, ethdkg, accusations) against the layer 1 blockchains (only ethereum at this point).
  * **Transaction**: responsible for watching for layer 1 transactions (only ethereum at this point) done by the AliceNet node and retrieve its receipts.

**Monitor:**

* Simplified the monitoring system.
* Removed all call to smart contracts during event processing.


#### Executor

* Unified ethdkg and snapshot tasks under a unique interface called **Task**.
* Task scheduler was decoupled from monitoring and now it's its own service under the executor sub-package.
* Task manager was completely refactored to work with the new Task interface and scheduler.
* Added persistence of state for the task manager and task scheduler. In case of a node exit, all running and future tasks will be saved and resume under node initialization.

#### Smart Contracts

* Added new features to the deployment scripts. Now it's possible to schedule a maintenance (required to change validators) and unregister validators during the local tests.


### Bug fixes:

Monitor
* Fixed race conditions around the event getter workers.
*
