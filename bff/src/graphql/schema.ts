
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
