package service

import (
	"context"
	"fmt"
	"io"

	database "github.com/marcodev1x/grpc-tests/internal"
	"github.com/marcodev1x/grpc-tests/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDb *database.Category
}

func NewCategoryService(categoryDb *database.Category) *CategoryService {
	return &CategoryService{CategoryDb: categoryDb}
}

func (c *CategoryService) CreateCategory(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDb.Create(request.Name, request.Description)

	fmt.Println(err)

	categoryResponse := &pb.CategoryResponse{
		Category: &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}

	return categoryResponse, nil
}

func (c *CategoryService) FindCategories(context.Context, *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDb.FindAll()
	if err != nil {
		return nil, err
	}

	categoriesList := make([]*pb.Category, 0, len(categories))
	for _, category := range categories {
		categoriesList = append(categoriesList, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryList{Categories: categoriesList}, nil
}

func (c *CategoryService) FindCategoryUnique(ctx context.Context, request *pb.FindCategoryUniqueRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDb.Find(request.Id)

	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Category: &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err == io.EOF { // Não tem mais nada pra mandar
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return stream.SendAndClose(categories)
		}

		categoryResult, err := c.CategoryDb.Create(category.Name, category.Description)

		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}
