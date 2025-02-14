# Storage-Users

Purpose and description to be added

## Deprecated Metadata Backend

Starting with ocis version 3.0.0, the default backend for metadata switched to messagepack. If the setting `STORAGE_USERS_OCIS_METADATA_BACKEND` has not been defined manually, the backend will be migrated to `messagepack` automatically. Though still possible to manually configure `xattrs`, this setting should not be used anymore as it will be removed in a later version.

## Graceful Shutdown

Starting with Infinite Scale version 3.1, you can define a graceful shutdown period for the `storage-users` service.

IMPORTANT: The graceful shutdown period is only applicable if the `storage-users` service runs as standalone service. It does not apply if the `storage-users` service runs as part of the single binary or as single Docker environment. To build an environment where the `storage-users` service runs as a standalone service, you must start two instances, one _without_ the `storage-users` service and one _only with_ the the `storage-users` service. Note that both instances must be able to communicate on the same network.

When hard-stopping Infinite Scale, for example with the `kill <pid>` command (SIGKILL), it is possible and likely that not all data from the decomposedfs (metadata) has been written to the storage which may result in an inconsistent decomposedfs. When gracefully shutting down Infinite Scale, using a command like SIGTERM, the process will no longer accept any write requests from _other_ services and will try to write the internal open  requests which can take an undefined duration based on many factors. To mitigate that situation, the following things have been implemented:

*   With the value of the environment variable `STORAGE_USERS_GRACEFUL_SHUTDOWN_TIMEOUT`, the `storage-users` service will delay its shutdown giving it time to finalize writing necessary data. This delay can be necessary if there is a lot of data to be saved and/or if storage access/thruput is slow. In such a case you would receive an error log entry informing you that not all data could be saved in time. To prevent such occurrences, you must increase the default value.

*   If a shutdown error has been logged, the command-line maintenance tool [Inspect and Manipulate Node Metadata](https://doc.owncloud.com/ocis/next/maintenance/commands/commands.html#inspect-and-manipulate-node-metadata) can help to fix the issue. Please contact support for details.

## CLI Commands

For any command listed, use `--help` to get more details and possible options and arguments.

To authenticate CLI commands use:

*   `OCIS_SERVICE_ACCOUNT_SECRET=<acc-secret>` and
*   `OCIS_SERVICE_ACCOUNT_ID=<acc-id>`.

The `storage-users` CLI tool uses the default address to establish the connection to the `gateway` service. If the connection fails, check your custom `gateway` service `GATEWAY_GRPC_ADDR` configuration and set the same address in `storage-users` `OCIS_GATEWAY_GRPC_ADDR` or `STORAGE_USERS_GATEWAY_GRPC_ADDR`.

### Manage Unfinished Uploads

<!-- referencing: [oCIS FS] clean up aborted uploads https://github.com/owncloud/ocis/issues/2622 -->

When using Infinite Scale as user storage, a directory named `storage/users/uploads` can be found in the Infinite Scale data folder. This is an intermediate directory based on [TUS](https://tus.io) which is an open protocol for resumable uploads. Each upload consists of a _blob_ and a _blob.info_ file. Note that the term _blob_ is just a placeholder.

*   **If an upload succeeds**, the _blob_ file will be moved to the target and the _blob.info_ file will be deleted.

*   **In case of incomplete uploads**, the _blob_ and _blob.info_ files will continue to receive data until either the upload succeeds in time or the upload expires based on the `STORAGE_USERS_UPLOAD_EXPIRATION` variable, see the table below for details.

*   **In case of expired uploads**, the _blob_ and _blob.info_ files will _not_ be removed automatically. Thus a lot of data can pile up over time wasting storage space.

*   **In the rare case of a failure**, after the upload succeeded but the file was not moved to its target location, which can happen when postprocessing fails, the situation is the same as with expired uploads.

Example cases for expired uploads

*   In the final step the upload blob is moved from the upload area to the final blobstore (e.g. S3). 

*   If the bandwidth is limited and the file to transfer can't be transferred completely before the upload expiration time is reached, the file expires and can't be processed.

The admin can restart the postprocessing for this with the postprocessing cli.

The storage users service can only list and clean upload sessions:

```bash
ocis storage-users uploads <command>
```

```plaintext
COMMANDS:
   sessions   Print a list of upload sessions
   clean      Clean up leftovers from expired uploads
   list       Print a list of all incomplete uploads (deprecated)
```

#### Command Examples

Command to list ongoing upload sessions

```bash
ocis storage-users sessions --expired=false
```

```plaintext
Not expired sessions:
+--------------------------------------+--------------------------------------+---------+--------+------+--------------------------------------+--------------------------------------+---------------------------+------------+
|                Space                 |              Upload Id               |  Name   | Offset | Size |              Executant               |                Owner                 |          Expires          | Processing |
+--------------------------------------+--------------------------------------+---------+--------+------+--------------------------------------+--------------------------------------+---------------------------+------------+
| f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | 5e387954-7313-4223-a904-bf996da6ec0b | foo.txt |      0 | 1234 | f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | 2024-01-26T13:04:31+01:00 | false      |
| f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | f066244d-97b2-48e7-a30d-b40fcb60cec6 | bar.txt |      0 | 4321 | f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c | 2024-01-26T13:18:47+01:00 | false      |
+--------------------------------------+--------------------------------------+---------+--------+------+--------------------------------------+--------------------------------------+---------------------------+------------+
```

The sessions command can also output json

```bash
ocis storage-users sessions --expired=false --json
```

```json
{"id":"5e387954-7313-4223-a904-bf996da6ec0b","space":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c","filename":"foo.txt","offset":0,"size":1234,"executant":{"idp":"https://cloud.ocis.test","opaque_id":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"},"spaceowner":{"opaque_id":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"},"expires":"2024-01-26T13:04:31+01:00","processing":false}
{"id":"f066244d-97b2-48e7-a30d-b40fcb60cec6","space":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c","filename":"bar.txt","offset":0,"size":4321,"executant":{"idp":"https://cloud.ocis.test","opaque_id":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"},"spaceowner":{"opaque_id":"f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"},"expires":"2024-01-26T13:18:47+01:00","processing":false}
```

Command to clear expired uploads
```bash
ocis storage-users uploads clean
```

```plaintext
Cleaned uploads:
- 455bd640-cd08-46e8-a5a0-9304908bd40a (Filename: file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

Deprecated list command to identify unfinished uploads

```bash
ocis storage-users uploads list
```

```plaintext
Incomplete uploads:
 - 455bd640-cd08-46e8-a5a0-9304908bd40a (file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

### Manage Trash-Bin Items

This command set provides commands to get an overview of trash-bin items, restore items and purge old items of `personal` spaces and `project` spaces (spaces that have been created manually). `trash-bin` commands require a `spaceID` as parameter. See [List all spaces
](https://owncloud.dev/apis/http/graph/spaces/#list-all-spaces-get-drives) or [Listing Space IDs](https://doc.owncloud.com/ocis/5.0/maintenance/space-ids/space-ids.html) for details of how to get them.

```bash
ocis storage-users trash-bin <command>
```

```plaintext
COMMANDS:
   purge-expired  Purge expired trash-bin items
   list           Print a list of all trash-bin items of a space.
   restore-all    Restore all trash-bin items for a space.
   restore        Restore a trash-bin item by ID.
```

#### Purge Expired

Purge all expired items from the trash-bin.

```bash
ocis storage-users trash-bin purge-expired
```

The behaviour of the `purge-expired` command can be configured by using the following environment variables.

*   `STORAGE_USERS_PURGE_TRASH_BIN_USER_ID`  
Used to obtain space trash-bin information and takes the system admin user as the default which is the `OCIS_ADMIN_USER_ID` but can be set individually. It should be noted, that the `OCIS_ADMIN_USER_ID` is only assigned automatically when using the single binary deployment and must be manually assigned in all other deployments. The command only considers spaces to which the assigned user has access and delete permission.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE`  
Has a default value of `720h` which equals `30 days`. This means, the command will delete all files older than `30 days`. The value is human-readable, for valid values see the duration type described in the [Environment Variable Types](https://doc.owncloud.com/ocis/5.0/deployment/services/envvar-types-description.html). A value of `0` is equivalent to disable and prevents the deletion of `personal space` trash-bin files.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE`  
Has a default value of `720h` which equals `30 days`. This means, the command will delete all files older than `30 days`. The value is human-readable, for valid values see the duration type described in the [Environment Variable Types](https://doc.owncloud.com/ocis/5.0/deployment/services/envvar-types-description.html). A value of `0` is equivalent to disable and prevents the deletion of `project space` trash-bin files.

#### List and Restore Trash-Bins Items

The variable `STORAGE_USERS_CLI_MAX_ATTEMPTS_RENAME_FILE` defines a maximum number of attempts to rename a file when the admin restores the file with the CLI option `--option keep-both` to an existing destination with the same name.

*   Print a list of all trash-bin items of a space
    ```bash
    ocis storage-users trash-bin list
    ```

* Restore all trash-bin items for a space
    ```bash
    ocis storage-users trash-bin restore-all
    ```

* Restore a trash-bin item by ID
    ```bash
    ocis storage-users trash-bin restore
    ```

## Caching

The `storage-users` service caches stat, metadata and uuids of files and folders via the configured store in `STORAGE_USERS_STAT_CACHE_STORE`, `STORAGE_USERS_FILEMETADATA_CACHE_STORE` and `STORAGE_USERS_ID_CACHE_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `redis-sentinel`: Stores data in a configured Redis Sentinel cluster.
  -   `nats-js-kv`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in production environments.
  -   `ocmem`: Advanced in-memory store allowing max size. (deprecated)
  -   `redis`: Stores data in a configured Redis cluster. (deprecated)
  -   `etcd`: Stores data in a configured etcd cluster. (deprecated)
  -   `nats-js`: Stores data using object-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/obj_store) (deprecated)

Other store types may work but are not supported currently.

Note: The service can only be scaled if not using `memory` store and the stores are configured identically over all instances!

Note that if you have used one of the deprecated stores, you should reconfigure to one of the supported ones as the deprecated stores will be removed in a later version.

Store specific notes:
  -   When using `redis-sentinel`, the Redis master to use is configured via e.g. `OCIS_CACHE_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.
  -   When using `nats-js-kv` it is recommended to set `OCIS_CACHE_STORE_NODES` to the same value as `OCIS_EVENTS_ENDPOINT`. That way the cache uses the same nats instance as the event bus.
  -   When using the `nats-js-kv` store, it is possible to set `OCIS_CACHE_DISABLE_PERSISTENCE` to instruct nats to not persist cache data on disc.
