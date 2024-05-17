package product

import (
	"eniqlo/internal/middleware"
	"eniqlo/pkg/response"
	"eniqlo/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	uc IProductUsecase
}

func NewProductHandler(uc IProductUsecase) *productHandler {
	return &productHandler{
		uc: uc,
	}
}

func (h *productHandler) Router(r *gin.RouterGroup) {
	group := r.Group("product")

	group.Use(middleware.UseJwtAuth)

	group.POST("", h.CreateProduct)
	group.GET("", h.FindAll)
	group.DELETE("/:id", h.DeleteProduct)
	group.GET("/customer", h.GetPublicProductHandler)
}

func (h *productHandler) FindAll(c *gin.Context) {
	query := QueryParams{}
	if err := c.ShouldBindQuery(&query); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	products, err := h.uc.FindProducts(query)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("Product fetched successfully!"), response.WithData(products))
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	product, err := h.uc.CreateProduct(request)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	res := FormatCreateProductResponse(*product)

	response.GenerateResponse(c, http.StatusCreated, response.WithMessage("Product created successfully!"), response.WithData(res))
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	userId := c.MustGet("userID").(string)

	if err := h.uc.DeleteProduct(id, userId); err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"))
}

func (h *productHandler) GetPublicProductHandler(c *gin.Context) {
	var queryParam ProductFilter

	limitFilter := c.DefaultQuery("limit", "5")
	offsetFilter := c.DefaultQuery("offset", "0")
	nameFilter := c.Query("name")
	categoryFilter := c.Query("category")
	skuFilter := c.Query("sku")
	priceFilter := c.Query("price")
	instockFilter := c.Query("inStock")

	queryParam.Limit = validation.ParseInt(limitFilter)
	queryParam.Offset = validation.ParseInt(offsetFilter)
	queryParam.Name = nameFilter
	queryParam.Category = ProductCategory(categoryFilter)
	queryParam.Sku = skuFilter
	queryParam.Price = SortBy(priceFilter)
	queryParam.InStock = InStockEnum(instockFilter)

	if err := c.ShouldBindQuery(&queryParam); err != nil {
		res := validation.FormatValidation(err)
		response.GenerateResponse(c, res.Code, response.WithMessage(res.Message))
		return
	}

	products, err := h.uc.GetPublicProducts(queryParam)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		return
	}

	productsResponse := []gin.H{}
	for _, product := range products {
		productsResponse = append(productsResponse, gin.H{
			"id":        product.ID,
			"name":      product.Name,
			"sku":       product.SKU,
			"category":  product.Category,
			"imageUrl":  product.ImageURL,
			"stock":     product.Stock,
			"price":     product.Price,
			"location":  product.Location,
			"createdAt": product.CreatedAt.Format("2006-01-02T15:04:05"),
		})
	}

	response.GenerateResponse(c, http.StatusOK, response.WithData(productsResponse))
}
