import {
  MigrateUpArgs,
  MigrateDownArgs,
} from "@payloadcms/db-mongodb";

export async function up({ payload }: MigrateUpArgs): Promise<void> {
};

export async function down({ payload }: MigrateDownArgs): Promise<void> {
  // Migration code
};
