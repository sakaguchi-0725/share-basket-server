
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

export class CreatePersonalShoppingItemInput {
    name: string;
    category: string;
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

export abstract class IMutation {
    abstract createPersonalShoppingItem(input: CreatePersonalShoppingItemInput): PersonalShoppingItem | Promise<PersonalShoppingItem>;
}

type Nullable<T> = T | null;
