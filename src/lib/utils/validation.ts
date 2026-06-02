
export type Sanitation<T> = {
    [K in keyof T]: T[K] extends { error: any } ? T[K]["error"] : never;
};

type Error = string | "";
type Value = string | number | any;

export type ValidationFunc = (value: any, ...args: any[]) => Error;

export type Fillable = {
    [key: string]: {
        value: Value;
        error: Error;
        validator: ValidationFunc;
    } | Value;
};

export function extract<T>(data: T): Sanitation<T> {
    const out: any = {};

    for (const key in data) {
        const entry = (data as any)[key];

        out[key] =
            entry && typeof entry === "object" && "value" in entry
                ? entry.value
                : entry;
    }

    return out;
}

export function validateValue(data: string, validator?: ValidationFunc): Error {

    if (!validator) return "";

    let error = validator(data);

    if (error.length > 0) {
        return error;
    } else {
        return "";
    }
}

export function validate(data: Fillable): Error {
    let elementWithError = Object.values(data).find(element =>
        [element.error, element.validator(element?.value)].some(thing => thing !== ''));

    if (elementWithError === undefined) {
        return "";
    }

    return elementWithError.error;
}

function isDigit(char: string) {
    return /^\d$/.test(char);
}

function filterDigits(str: string): string {
    return str.split('').filter(e => isDigit(e)).join('');
}

function sanitizeString(str: string): string {
    return str.replace(/[^A-Za-z]/g, "");
}

function isSanitized(str: string): boolean {
    return str === sanitizeString(str);
}

function isValidDate(x: string): boolean {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(x)) return false;

    const date = new Date(x + "T00:00:00");
    if (isNaN(date.getTime())) return false;

    const [year, month, day] = x.split("-").map(Number);
    return (
        date.getUTCFullYear() === year &&
        date.getUTCMonth() + 1 === month &&
        date.getUTCDate() === day
    );
}

function isValidPhoneNumber(phone: string) {

    phone = phone.split(' ').join('');

    if (!phone.startsWith('0')) {
        return "Your phone number must start with 0."
    }

    if (phone.split('').some(character => !isDigit(character))) {
        return "Phone number must only contain digits";
    }

    if (phone.length !== 10) {
        return "Phone number must be in the format 0XXX XXX XXX";
    }

    return "";

}
function isValidPrescription(index: number, value: string): Error {
    if (value.trim() == "") {
        return "Prescription must not be empty."
    }
    return "";
}

function parseISODate(dateStr: string): Date | null {
    if (!dateStr) return null;

    const parts = dateStr.split("-");
    if (parts.length !== 3) return null;

    const [yearStr, monthStr, dayStr] = parts;
    const year = Number(yearStr);
    const month = Number(monthStr);
    const day = Number(dayStr);

    if (!year || !month || !day) return null;

    const d = new Date(year, month - 1, day);
    d.setHours(0, 0, 0, 0);


    if (d.getFullYear() !== year || d.getMonth() !== month - 1 || d.getDate() !== day) {
        return null;
    }

    return d;
}

export const validation: Record<string, ValidationFunc> = {

    notEmpty(string: string): Error {
        return string.trim() !== "" ? "" : "This field cannot be empty.";
    },

    nothing(_: string): Error {
        return "";
    },

    email(email: string): Error {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email) ? "" : "Please enter a valid email.";
    },

    password(pwd: string): Error {
        return pwd.length > 6 ? "" : "Password must be greater than 6.";
    },

    age(date: string): Error {
        if (!isValidDate(date)) return "Please enter a valid date."

        let d = new Date(date);

        if (d.getDate() > new Date().getDate()) {
            return "Cannot accept future dates."
        }

        return "";
    },

    number(value: string): Error {
        const num = Number(value);
        if (isNaN(num)) {
            return "Please enter a valid number.";
        }
        if (num < 0) {
            return "Number cannot be negative.";
        }
        return "";
    },

    fullName(name: string): Error {

        if (!isSanitized(name)) return "Invalid characters found.";


        let names = name.split(' ');
        if (names.length < 2) return "Please enter your full name.";
        if (names.some(name => name.length < 3)) return "Name is too short."

        return "";
    },

    name(name: string): Error {

        if (name.length < 3) {
            return "Name is too short."
        }

        return "";

    },

    phoneNumber(phone: string): Error {
        return isValidPhoneNumber(phone);
    },

    futureDate(date: string): string {
        const selected = parseISODate(date);
        if (!selected) return "Please enter a valid date.";

        const today = new Date();
        today.setHours(0, 0, 0, 0);

        if (selected < today) {
            return "Booking date cannot be in the past.";
        }
        return "";
    },

    pastDate(date: string): string {
        const selected = parseISODate(date);
        if (!selected) return "Please enter a valid date.";

        const today = new Date();
        today.setHours(0, 0, 0, 0);

        if (selected > today) {
            return "Date cannot be in the future.";
        }
        return "";
    },


    isValidPrescription,

    studentNumber(value: string): Error {
        if (!/^\d{8,12}$/.test(value.trim())) {
            return "Student number must be 8-12 digits.";
        }
        return "";
    },

    pfeCode(value: string): Error {
        if (!/^PFE-\d{2}-[A-Z]+-\d{4}$/.test(value.trim())) {
            return "PFE code must follow format PFE-XX-SPEC-YEAR.";
        }
        return "";
    },

    algerianPhone(phone: string): Error {
        const cleaned = phone.replace(/\s/g, '');
        if (!/^(0|\+213)[5-7]\d{8}$/.test(cleaned)) {
            return "Enter a valid Algerian phone number (0XXXXXXXXX or +213XXXXXXXXX).";
        }
        return "";
    },

    universityEmail(email: string): Error {
        const emailErr = validation.email(email);
        if (emailErr) return emailErr;
        if (!email.endsWith('.dz')) {
            return "Must be a .dz university email address.";
        }
        return "";
    },

    grade(value: string): Error {
        const num = Number(value);
        if (isNaN(num)) return "Please enter a valid number.";
        if (num < 0 || num > 20) return "Grade must be between 0 and 20.";
        return "";
    },

    url(value: string): Error {
        if (!value.trim()) return "";
        try {
            new URL(value);
            return "";
        } catch {
            return "Please enter a valid URL.";
        }
    },
}
