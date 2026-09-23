export { M1Client } from "./client.js";
export { AbdmError } from "./common/errors.js";
export type { CallOptions, Environment, M1Config } from "./common/config.js";

export { AadhaarService } from "./aadhaar/aadhaar.service.js";
export { AddressService } from "./address/address.service.js";
export { LoginService } from "./login/login.service.js";
export { PhrService } from "./phr/phr.service.js";
export { ProfileService } from "./profile/profile.service.js";

export type * from "./common/types.js";
export type * from "./aadhaar/aadhaar.types.js";
export type * from "./address/address.types.js";
export type * from "./login/login.types.js";
export type * from "./phr/phr.types.js";
export type * from "./profile/profile.types.js";
