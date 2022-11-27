declare type bool = boolean;
declare type i8 = number;
declare type i16 = number;
declare type i32 = number;
declare type i64 = bigint;
declare type u8 = number;
declare type u16 = number;
declare type u32 = number;
declare type u64 = bigint;
declare type usize = number;


/** Converts any other numeric value to an 8-bit signed integer. */
declare function i8(value: any): i8;
declare namespace i8 {
    /** Smallest representable value. */
    export const MIN_VALUE: i8;
    /** Largest representable value. */
    export const MAX_VALUE: i8;
}

/** Converts any other numeric value to an 16-bit signed integer. */
declare function i16(value: any): i16;
declare namespace i16 {
    /** Smallest representable value. */
    export const MIN_VALUE: i16;
    /** Largest representable value. */
    export const MAX_VALUE: i16;
}

/** Converts any other numeric value to an 32-bit signed integer. */
declare function i32(value: any): i32;
declare namespace i32 {
    /** Smallest representable value. */
    export const MIN_VALUE: i32;
    /** Largest representable value. */
    export const MAX_VALUE: i32;
}

/** Converts any other numeric value to an 64-bit signed integer. */
declare function i64(value: any): i64;
declare namespace i64 {
    /** Smallest representable value. */
    export const MIN_VALUE: i64;
    /** Largest representable value. */
    export const MAX_VALUE: i64;
}

/** Converts any other numeric value to an 8-bit unsigned integer. */
declare function u8(value: any): u8;
declare namespace u8 {
    /** Smallest representable value. */
    export const MIN_VALUE: u8;
    /** Largest representable value. */
    export const MAX_VALUE: u8;
}

/** Converts any other numeric value to an 16-bit unsigned integer. */
declare function u16(value: any): u16;
declare namespace u16 {
    /** Smallest representable value. */
    export const MIN_VALUE: u16;
    /** Largest representable value. */
    export const MAX_VALUE: u16;
}

/** Converts any other numeric value to an 32-bit unsigned integer. */
declare function u32(value: any): u32;
declare namespace u32 {
    /** Smallest representable value. */
    export const MIN_VALUE: u32;
    /** Largest representable value. */
    export const MAX_VALUE: u32;
}

/** Converts any other numeric value to an 64-bit unsigned integer. */
declare function u64(value: any): u64;
declare namespace u64 {
    /** Smallest representable value. */
    export const MIN_VALUE: u64;
    /** Largest representable value. */
    export const MAX_VALUE: u64;
}
