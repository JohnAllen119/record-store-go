#include <stdio.h>
#include <stdlib.h>
typedef struct Node
{
    int data;
    struct Node *next;
} Node;
Node *createNode(int value)
{
    Node *newNode = malloc(sizeof(Node));
    if (newNode == NULL)
        return NULL;
    newNode->data = value;
    newNode->next = NULL;
    return newNode;
}
void freeList(Node *head)
{
    Node *current = head;
    while (current != NULL)
    {
        Node *nextNode = current->next;
        free(current);
        current = nextNode;
    }
}
Node *insertAtHead(Node *head, int value)
{
    Node *newNode = createNode(value);
    if (newNode == NULL)
        return head;
    newNode->next = head;
    return newNode;
}
Node *insertAtTail(Node *head, int value)
{
    Node *newNode = createNode(value);
    if (newNode == NULL)
        return head;
    if (head == NULL)
    {
        return newNode;
    }
    Node *current = head;
    while (current->next)
    {
        current = current->next;
    }
    current->next = newNode;
    return head;
}
Node *findNode(Node *head, int value)
{
    Node *current = head;
    while (current != NULL)
    {
        if (current->data == value)
        {
            return current;
        }
        current = current->next;
    }
    return NULL;
}
Node *deleteHead(Node *head)
{
    if (head == NULL)
    {
        return NULL;
    }
    Node *deletedNode = head;
    head = head->next;

    free(deletedNode);
    return head;
}
Node *deleteByValue(Node *head, int value){
    if(head==NULL){
        return NULL;
    }
    if(head->data ==value){
        return deleteHead(head);
    }
    Node *current=head;
    while(current->next != NULL && current->next->data != value){
        current=current->next;
    }
    if(current->next !=NULL){
        Node* deleteNode = current->next;
        current->next=deleteNode->next;
        free(deleteNode);
    }
    return head;
}
int main(void)
{
    Node *head = createNode(10);
    Node *second = createNode(20);
    Node *third = createNode(30);

    if (head == NULL || second == NULL || third == NULL)
    {
        free(head);
        free(second);
        free(third);
        return 1;
    }
    head->next = second;
    second->next = third;
    head = insertAtHead(head, 5);
    head = insertAtTail(head, 40);
    /* Node *current = head;

     while (current)
     {
         printf("%d\n", current->data);
         current = current->next;
     }

     Node *found = findNode(head, 99);

     if (found != NULL)
     {
         printf("找到节点:%d\n", found->data);
     }
     else
     {
         printf("没有找到\n");
     }*/
    head = deleteHead(head);
    head = deleteByValue(head, 20);
    Node *current = head;

    while (current)
    {
        printf("%d\n", current->data);
        current = current->next;
    }

    freeList(head);
    head = NULL;
    return 0;
}