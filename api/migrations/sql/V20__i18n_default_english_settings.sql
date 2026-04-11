WITH cleared_defaults AS (
  UPDATE languages
  SET is_default = FALSE,
      updated_at = NOW()
  WHERE deleted = FALSE
    AND code <> 'en'
    AND is_default = TRUE
),
upserted_language AS (
  INSERT INTO languages (
    code,
    name,
    native_name,
    is_default,
    active,
    deleted,
    created_at,
    updated_at
  )
  VALUES (
    'en',
    'English',
    'English',
    TRUE,
    TRUE,
    FALSE,
    NOW(),
    NOW()
  )
  ON CONFLICT (code, deleted) DO UPDATE
  SET name = EXCLUDED.name,
      native_name = EXCLUDED.native_name,
      is_default = EXCLUDED.is_default,
      active = EXCLUDED.active,
      deleted = FALSE,
      updated_at = NOW()
  RETURNING id
)
SELECT id FROM upserted_language;
