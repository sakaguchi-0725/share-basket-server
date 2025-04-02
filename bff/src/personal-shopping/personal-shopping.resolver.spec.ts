import { Test, TestingModule } from '@nestjs/testing';
import { PersonalShoppingResolver } from './personal-shopping.resolver';
import { PersonalShoppingService } from './personal-shopping.service';
import {
  CreatePersonalShoppingItemInput,
  PersonalShoppingItem,
  PersonalShoppingStatus,
} from 'src/graphql/schema';
import { AuthService } from 'src/auth/auth.service';
import { AuthGuard } from 'src/auth/auth.guard';

describe('PersonalShoppingResolver', () => {
  let resolver: PersonalShoppingResolver;
  let service: PersonalShoppingService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        PersonalShoppingResolver,
        {
          provide: PersonalShoppingService,
          useValue: {
            findAll: jest.fn(),
            create: jest.fn(),
          },
        },
        {
          provide: AuthService,
          useValue: {
            verifyToken: jest.fn().mockReturnValue('mock-user-id'),
          },
        },
        {
          provide: AuthGuard,
          useValue: {
            canActivate: jest.fn().mockReturnValue(true),
          },
        },
      ],
    }).compile();

    resolver = module.get<PersonalShoppingResolver>(PersonalShoppingResolver);
    service = module.get<PersonalShoppingService>(PersonalShoppingService);
  });

  it('個人用の買い物リストを全て取得できること', () => {
    const mockData: PersonalShoppingItem[] = [
      {
        id: 1,
        name: '牛乳',
        category: 'food',
        status: PersonalShoppingStatus.UNPURCHASED,
      },
      {
        id: 2,
        name: 'トイレットペーパー',
        category: 'daily',
        status: PersonalShoppingStatus.PURCHASED,
      },
    ];

    jest.spyOn(service, 'findAll').mockReturnValue(mockData);

    const reuslt = resolver.getPersonalShoppingItems();
    expect(service.findAll).toHaveBeenCalled();
    expect(reuslt).toEqual(mockData);
  });

  it('個人用の買い物リストに新しいアイテムを追加できること', () => {
    const input: CreatePersonalShoppingItemInput = {
      name: '牛乳',
      category: 'foods',
    };

    const mockData: PersonalShoppingItem = {
      id: 6,
      name: '牛乳',
      status: PersonalShoppingStatus.UNPURCHASED,
      category: 'foods',
    };

    jest.spyOn(service, 'create').mockReturnValue(mockData);

    const result = resolver.createPersonalShoppingItem(input);
    expect(service.create).toHaveBeenCalledWith(input);
    expect(result).toEqual(mockData);
  });
});
