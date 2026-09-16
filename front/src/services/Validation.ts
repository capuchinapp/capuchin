import { get } from "svelte/store";
import * as yup from "yup";
import { t } from "../i18n";

type ValidateCallback = () => void;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type ValidateErrorCallback = (errors: any) => void;

export default class Validation {
  #msgRequiredInput = "";
  #msgRequiredSelect = "";

  #loginSchema?: yup.AnySchema;
  #registerSchema?: yup.AnySchema;
  #applySchema?: yup.AnySchema;
  #clientSchema?: yup.AnySchema;
  #projectSchema?: yup.AnySchema;
  #timelogSchema?: yup.AnySchema;

  #requiredEmail?: yup.StringSchema;
  #requiredName?: yup.StringSchema;
  #requiredSelectId?: yup.StringSchema;
  #requiredBillableRate?: yup.NumberSchema;

  #ulidRegex = /^[0-9A-HJKMNP-TV-Z]{26}$/;

  constructor() {
    this.#msgRequiredInput = get(t)("validate.required.input");
    this.#msgRequiredSelect = get(t)("validate.required.select");
  }

  validateLogin(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#loginSchema) {
      this.#loginSchema = yup.object().shape({
        email: this.#getRequiredEmail(),
      });
    }

    this.#validateSchema(
      this.#loginSchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  validateRegister(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#registerSchema) {
      this.#registerSchema = yup.object().shape({
        email: this.#getRequiredEmail(),
      });
    }

    this.#validateSchema(
      this.#registerSchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  validateApplyCode(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#applySchema) {
      this.#applySchema = yup.object().shape({
        email: this.#getRequiredEmail(),
        code: yup.string().required(this.#msgRequiredInput),
      });
    }

    this.#validateSchema(
      this.#applySchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  validateClient(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#clientSchema) {
      this.#clientSchema = yup.object().shape({
        name: this.#getRequiredName(),
        billableRate: this.#getRequiredBillableRate(),
      });
    }

    this.#validateSchema(
      this.#clientSchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  validateProject(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#projectSchema) {
      this.#projectSchema = yup.object().shape({
        clientId: this.#getRequiredSelectId(),
        name: this.#getRequiredName(),
        billableRate: this.#getRequiredBillableRate(),
      });
    }

    this.#validateSchema(
      this.#projectSchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  validateTimelog(
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    if (!this.#timelogSchema) {
      const time = yup
        .string()
        .matches(/\d{0,23}:\d{0,59}/, { message: get(t)("validate.time") });

      this.#timelogSchema = yup.object().shape({
        projectId: this.#getRequiredSelectId(),
        date: yup.date().required(this.#msgRequiredInput),
        timeStart: time.required(this.#msgRequiredInput),
        timeEnd: time.nullable(),
        billableRate: this.#getRequiredBillableRate(),
      });
    }

    this.#validateSchema(
      this.#timelogSchema,
      data,
      callbackSuccess,
      callbackError,
    );
  }

  #validateSchema(
    scheme: yup.AnySchema,
    data: Record<string, unknown>,
    callbackSuccess: ValidateCallback,
    callbackError: ValidateErrorCallback,
  ): void {
    scheme
      .validate(data, { abortEarly: false })
      .then(() => callbackSuccess())
      .catch((err) => callbackError(this.#extractErrors(err)));
  }

  #getRequiredEmail(): yup.StringSchema {
    if (!this.#requiredEmail) {
      this.#requiredEmail = yup
        .string()
        .required(this.#msgRequiredInput)
        .email(get(t)("validate.email"));
    }

    return this.#requiredEmail;
  }

  #getRequiredName(): yup.StringSchema {
    if (!this.#requiredName) {
      this.#requiredName = yup
        .string()
        .required(this.#msgRequiredInput)
        .max(255, get(t)("validate.string.max255"));
    }

    return this.#requiredName;
  }

  #getRequiredSelectId(): yup.StringSchema {
    if (!this.#requiredSelectId) {
      this.#requiredSelectId = yup
        .string()
        .matches(this.#ulidRegex, get(t)("validate.ulid"))
        .required(this.#msgRequiredSelect);
    }

    return this.#requiredSelectId;
  }

  #getRequiredBillableRate(): yup.NumberSchema {
    if (!this.#requiredBillableRate) {
      this.#requiredBillableRate = yup
        .number()
        .required(this.#msgRequiredInput)
        .min(0, get(t)("validate.number.gtezero"));
    }

    return this.#requiredBillableRate;
  }

  #extractErrors(err: yup.ValidationError): Record<string, string> {
    return err.inner.reduce(
      (acc, cur: yup.ValidationError) => {
        return { ...acc, [cur.path as string]: cur.message };
      },
      {} as Record<string, string>,
    );
  }
}
