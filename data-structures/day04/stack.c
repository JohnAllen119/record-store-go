#include <stdio.h>
#include <stdlib.h>

typedef struct Node
{
    int data;
    struct Node *next;
} Node;
Node *createNode(int data)
{
    Node *newNode = malloc(sizeof(Node));
    if (newNode == NULL)
        return NULL;
    newNode->data = data;
    newNode->next = NULL;
    return newNode;
}
Node *push(Node *top, int data)
{
    Node *newNode = createNode(data);
    if (newNode == NULL)
        return top;

    newNode->next = top;
    return newNode;
}
int peek(Node *top, int *value)
{
    if (top == NULL)
        return 0;
    *value = top->data;
    return 1;
}
int pop(Node **top, int *value)
{
    if (top == NULL || *top == NULL)
        return 0;
    Node *tmpNode = *top;
    *value = tmpNode->data;
    *top = tmpNode->next;
    free(tmpNode);
    return 1;
}
int isEmpty(Node *top)
{   
    
    return top == NULL;
}
int main(void)
{
    Node *top = NULL;
    top = push(top, 10);
    top = push(top, 20);
    top = push(top, 30);
    Node *current = top;
    while (current != NULL)
    {
        printf("%d\n", current->data);
        current = current->next;
    }
    int value;

    if (peek(top, &value))
    {
        printf("当前栈顶：%d\n", value);
    }

    while (pop(&top, &value))
    {
        printf("出栈：%d\n", value);
    }

    if (top == NULL)
    {
        printf("栈已清空\n");
    }
    return 0;
}