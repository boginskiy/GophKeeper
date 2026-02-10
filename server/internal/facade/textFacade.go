package facade

import (
	"context"

	"github.com/boginskiy/GophKeeper/server/internal/codec"
	"github.com/boginskiy/GophKeeper/server/internal/errs"
	"github.com/boginskiy/GophKeeper/server/internal/handler"
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
)

type TextFacade struct {
	rpc.UnimplementedTexterServiceServer
	ErrMapper   errs.ErrMapper
	TextCoder   codec.TextGRPCCoder[*model.Text]
	TextHandler handler.TextHandler[*model.Text]
}

func NewTextFacade(
	errMapper errs.ErrMapper,
	textCoder codec.TextGRPCCoder[*model.Text],
	textHandler handler.TextHandler[*model.Text],
) *TextFacade {

	return &TextFacade{
		ErrMapper:   errMapper,
		TextCoder:   textCoder,
		TextHandler: textHandler,
	}
}

func (t *TextFacade) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	reqModText, _ := t.TextCoder.DecoderCreateRequest(req)
	resModText, err := t.TextHandler.HandleCreate(ctx, reqModText)
	if err != nil {
		return nil, t.ErrMapper.Mapping(err)
	}

	return t.TextCoder.EncodeCreateResponse(resModText)
}

func (t *TextFacade) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	reqModText, _ := t.TextCoder.DecoderReadRequest(req)

	resModText, err := t.TextHandler.HandleRead(ctx, reqModText)
	if err != nil {
		return nil, t.ErrMapper.Mapping(err)
	}

	return t.TextCoder.EncodeReadResponse(resModText)
}

func (t *TextFacade) ReadAll(ctx context.Context, req *rpc.ReadAllRequest) (*rpc.ReadAllResponse, error) {
	reqModText, _ := t.TextCoder.DecoderReadAllRequest(req)

	resModTexts, err := t.TextHandler.HandleReadAll(ctx, reqModText)
	if err != nil {
		return nil, t.ErrMapper.Mapping(err)
	}

	return t.TextCoder.EncodeReadAllResponse(resModTexts)
}

func (t *TextFacade) Update(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	reqModText, _ := t.TextCoder.DecoderCreateRequest(req)
	resModText, err := t.TextHandler.HandleUpdate(ctx, reqModText)
	if err != nil {
		return nil, t.ErrMapper.Mapping(err)
	}

	return t.TextCoder.EncodeCreateResponse(resModText)
}

func (t *TextFacade) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.DeleteResponse, error) {
	reqModText, _ := t.TextCoder.DecoderDeleteRequest(req)

	resModText, err := t.TextHandler.HandleDelete(ctx, reqModText)
	if err != nil {
		return nil, t.ErrMapper.Mapping(err)
	}

	return t.TextCoder.EncodeDeleteResponse(resModText)
}
