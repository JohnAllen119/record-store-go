#include <stdio.h>
#define MAX_SIZE 100
typedef struct
{
    int data[MAX_SIZE];
    int top;
} SeqStack;

void initStack(SeqStack *stack)
{
    stack->top = -1;
}
int isEmpty(SeqStack *stack)
{
    return stack->top == -1;
}
int isFull(SeqStack *stack)
{
    return stack->top == MAX_SIZE - 1;
}
int push(SeqStack *stack, int value)
{
    if (isFull(stack))
        return 0;
    stack->top += 1;
    stack->data[stack->top] = value;
    return 1;
}
int peek(SeqStack *stack, int *value)
{
    if (isEmpty(stack))
        return 0;
    *value = stack->data[stack->top];
    return 1;
}
int pop(SeqStack *stack, int *value)
{
    if (isEmpty(stack))
        return 0;
    *value = stack->data[stack->top];
    stack->top--;
    return 1;
}
void swap(int *a,int *b){
    int c=0;
    c=*a;
    *a=*b;
    *b=c;
}
int main(void)
{
    SeqStack stack;
    initStack(&stack);
    push(&stack, 10);
    push(&stack, 20);
    push(&stack, 30);
    int value;
    if (peek(&stack, &value))
    {
        printf("peek: %d\n", value);
    }

    while (pop(&stack, &value))
    {
        printf("pop: %d\n", value);
    }

    printf("empty: %d\n", isEmpty(&stack));
    printf("pop empty: %d\n", pop(&stack, &value));

    return 0;
}