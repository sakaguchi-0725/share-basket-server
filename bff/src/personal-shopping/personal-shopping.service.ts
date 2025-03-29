import { Injectable } from '@nestjs/common';
import {
  CreatePersonalShoppingItemInput,
  PersonalShoppingItem,
  PersonalShoppingStatus,
} from 'src/graphql/schema';

@Injectable()
export class PersonalShoppingService {
  findAll(this: void): PersonalShoppingItem[] {
    return [
      {
        id: 1,
        name: '牛乳',
        status: PersonalShoppingStatus.UNPURCHASED,
        category: 'foods',
      },
      {
        id: 2,
        name: 'パン',
        status: PersonalShoppingStatus.UNPURCHASED,
        category: 'foods',
      },
      {
        id: 3,
        name: 'トイレットペーパー',
        status: PersonalShoppingStatus.UNPURCHASED,
        category: 'daily',
      },
      {
        id: 4,
        name: '食器用洗剤',
        status: PersonalShoppingStatus.UNPURCHASED,
        category: 'daily',
      },
      {
        id: 5,
        name: 'ヒートテック',
        status: PersonalShoppingStatus.UNPURCHASED,
        category: 'clothes',
      },
    ];
  }

  create(
    this: void,
    input: CreatePersonalShoppingItemInput,
  ): PersonalShoppingItem {
    return {
      id: 6,
      name: input.name,
      status: PersonalShoppingStatus.UNPURCHASED,
      category: input.category,
    };
  }
}
