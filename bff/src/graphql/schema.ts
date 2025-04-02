
/*
 * -------------------------------------------------------
 * THIS FILE WAS AUTOMATICALLY GENERATED (DO NOT MODIFY)
 * -------------------------------------------------------
 */

/* tslint:disable */
/* eslint-disable */

export enum PersonalShoppingStatus {
    UNPURCHASED = "UNPURCHASED",
    PURCHASED = "PURCHASED"
}

export class LoginInput {
    email: string;
    password: string;
}

export class SignUpInput {
    name: string;
    email: string;
    password: string;
}

export class CreatePersonalShoppingItemInput {
    name: string;
    category: string;
}

export abstract class IMutation {
    abstract signUp(input?: Nullable<SignUpInput>): boolean | Promise<boolean>;

    abstract login(input: LoginInput): boolean | Promise<boolean>;

    abstract createPersonalShoppingItem(input: CreatePersonalShoppingItemInput): PersonalShoppingItem | Promise<PersonalShoppingItem>;
}

export class PersonalShoppingItem {
    id: number;
    name: string;
    status: PersonalShoppingStatus;
    category: string;
}

export abstract class IQuery {
    abstract getPersonalShoppingItems(): PersonalShoppingItem[] | Promise<PersonalShoppingItem[]>;
}

type Nullable<T> = T | null;
