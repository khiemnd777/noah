CREATE INDEX IF NOT EXISTS ix_department_slug_not_deleted
  ON departments(slug)
  WHERE deleted = FALSE;
