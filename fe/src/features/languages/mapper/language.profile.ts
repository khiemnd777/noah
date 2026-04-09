import { mapper } from "@core/mapper/auto-mapper";
import type { LanguageModel, LanguageResourceModel } from "@features/languages/model/language.model";

mapper.register<LanguageResourceModel>({
  name: "LanguageResource",
  dtoToModelNaming: "snake_to_camel",
  modelToDtoNaming: "camel_to_snake",
  defaultModel() {
    return {
      id: 0,
      key: "",
      value: "",
      createdAt: null,
      updatedAt: null,
    };
  },
});

mapper.register<LanguageModel>({
  name: "Language",
  dtoToModelNaming: "snake_to_camel",
  modelToDtoNaming: "camel_to_snake",
  defaultModel() {
    return {
      id: 0,
      code: "",
      name: "",
      nativeName: "",
      isDefault: false,
      active: true,
      createdAt: null,
      updatedAt: null,
      resources: [],
    };
  },
});
