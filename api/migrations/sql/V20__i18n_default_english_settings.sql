WITH cleared_defaults AS (
  UPDATE languages
  SET is_default = FALSE,
      updated_at = NOW()
  WHERE deleted = FALSE
    AND code <> 'en'
    AND is_default = TRUE
),
updated_language AS (
  UPDATE languages
  SET name = 'English',
      native_name = 'English',
      is_default = TRUE,
      active = TRUE,
      deleted = FALSE,
      updated_at = NOW()
  WHERE code = 'en'
    AND deleted = FALSE
  RETURNING id
),
inserted_language AS (
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
  SELECT
    'en',
    'English',
    'English',
    TRUE,
    TRUE,
    FALSE,
    NOW(),
    NOW()
  WHERE NOT EXISTS (SELECT 1 FROM updated_language)
  RETURNING id
)
SELECT id FROM updated_language
UNION ALL
SELECT id FROM inserted_language;
